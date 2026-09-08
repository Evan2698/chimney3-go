package kcpproxy

import (
	"chimney3-go/core"
	servercontext "chimney3-go/sesrvercontext"
	"chimney3-go/settings"
	"chimney3-go/utils"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/xtaci/kcp-go/v5"
)

func runKCPClientCtx(ctx context.Context, s settings.Settings) error {

	httpAddr := s.Httpurl
	socks5Url := fmt.Sprintf("socks5://%s", s.Listen)
	log.Printf("Starting HTTP to SOCKS5 proxy on %s forwarding to %s", httpAddr, socks5Url)
	go core.Run2HTTP(httpAddr, socks5Url)

	listenAddress := s.Listen

	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Printf("Error listening on %s: %v", listenAddress, err)
		return err
	}
	// ensure listener is closed when context is done
	go func() {
		<-ctx.Done()
		l.Close()
	}()
	defer l.Close()
	log.Printf("KCP client listening on %s", listenAddress)

	key := deriveKey(s.Username)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		log.Printf("Error creating block cipher: %v", err)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				// listener closed due to context cancellation
				break
			default:
				log.Printf("Error accepting connection: %v", err)
			}
			break
		}
		go handleKCPConnection(conn, s, block)
	}

	return nil
}

func handleKCPConnection(conn net.Conn, s settings.Settings, block kcp.BlockCrypt) {

	defer conn.Close()
	proxyAddr := s.RemoteListen
	log.Printf("Connecting to KCP server at %s", proxyAddr)

	sess, err := kcp.DialWithOptions(proxyAddr, block, 10, 3)
	if err != nil {
		log.Printf("Error dialing KCP server: %v", err)
		return
	}
	defer sess.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go clientRoutine(conn, sess, &wg)
	go clientRoutine(sess, conn, &wg)
	wg.Wait()
}

func clientRoutine(src, dest net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer src.Close()
	defer dest.Close()
	io.Copy(dest, src)
}

func runKCPServerCtx(ctx context.Context, s settings.Settings) error {
	listenAddress := s.Listen

	key := deriveKey(s.Username)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		log.Printf("Error creating block cipher: %v", err)
		return err
	}

	l, err := kcp.ListenWithOptions(listenAddress, block, 10, 3)
	if err != nil {
		log.Printf("Error listening on %s: %v", listenAddress, err)
		return err
	}
	// ensure listener is closed when context is done
	go func() {
		<-ctx.Done()
		l.Close()
	}()
	defer l.Close()
	log.Printf("KCP server listening on %s", listenAddress)

	utils.StartUDPServerIfConfigured(s.Udplisten)

	for {
		sess, err := l.AcceptKCP()
		if err != nil {
			select {
			case <-ctx.Done():
				// listener closed due to context cancellation
				break
			default:
				log.Printf("Error accepting KCP connection: %v", err)
			}
			break
		}
		go handleKCPServerSession(sess)
	}

	return nil
}

// handleKCPServerSession handles an incoming KCP session by echoing data back to the client.
func handleKCPServerSession(conn *kcp.UDPSession) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(" fatal error on udp server: ", err)
		}
	}()

	defer conn.Close()
	// 1. 握手阶段
	header := make([]byte, 2)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return
	}
	// 检查 SOCKS5 版本
	if header[0] != 0x05 {
		return
	}
	// 不认证
	conn.Write([]byte{0x05, 0x00})

	// 2. 请求阶段
	requestHeader := make([]byte, 4)
	if _, err = io.ReadFull(conn, requestHeader); err != nil {
		return
	}
	if requestHeader[0] != 0x05 || requestHeader[1] != 0x01 { // 只支持 CONNECT
		return
	}
	var addr string
	switch requestHeader[3] {
	case 0x01: // IPv4
		address := make([]byte, 6)
		if _, err = io.ReadFull(conn, address); err != nil {
			return
		}
		addr = net.IP(address[:4]).String()
		port := binary.BigEndian.Uint16(address[4:6])
		addr = net.JoinHostPort(addr, fmt.Sprint(int(port)))
	case 0x03: // 域名
		length := []byte{0}
		if _, err = io.ReadFull(conn, length); err != nil {
			return
		}
		domainAndPort := make([]byte, int(length[0])+2)
		if _, err = io.ReadFull(conn, domainAndPort); err != nil {
			return
		}
		addr = string(domainAndPort[:length[0]])
		port := binary.BigEndian.Uint16(domainAndPort[length[0]:])
		addr = net.JoinHostPort(addr, fmt.Sprint(int(port)))
	case 0x04: // IPv6
		address := make([]byte, 18)
		if _, err = io.ReadFull(conn, address); err != nil {
			return
		}
		addr = net.IP(address[:16]).String()
		port := binary.BigEndian.Uint16(address[16:])
		addr = net.JoinHostPort(addr, fmt.Sprint(int(port)))
	default:
		return
	}

	// 3. 连接目标服务器
	target, err := net.Dial("tcp", addr)
	if err != nil {
		// 连接失败
		conn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}
	defer target.Close()
	// 连接成功
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})

	var wg sync.WaitGroup
	wg.Add(2)
	go clientRoutine(conn, target, &wg)
	go clientRoutine(target, conn, &wg)
	wg.Wait()
}

func RunKCPRoutine(s *settings.Settings, ctx servercontext.ServerContext) error {
	isServer := utils.IsServerMode(s.Mode)
	if isServer {
		return runKCPServerCtx(ctx, *s)
	}
	return runKCPClientCtx(ctx, *s)
}
