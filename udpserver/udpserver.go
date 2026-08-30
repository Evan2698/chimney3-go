package udpserver

import (
	"context"
	"log"
	"net"
	"sync/atomic"
	"time"
	"tun2proxylib/gvisorcore/buffer"
	"tun2proxylib/udppackage"
)

var (
	stop int32
)

const (
	defaultUDPURL = "0.0.0.0:5353"
	timeout       = 20 // seconds
)

func resetStopFlag() {
	atomic.StoreInt32(&stop, 0)
}

func resolveUDPAddress(udpURL string) (*net.UDPAddr, error) {
	if udpURL == "" {
		udpURL = defaultUDPURL
	}
	return net.ResolveUDPAddr("udp", udpURL)
}

func serveLoop(ctx context.Context, conn *net.UDPConn) {
	buf := buffer.Get()
	defer buffer.Put(buf)

	for {
		if atomic.LoadInt32(&stop) != 0 {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		target, src, payload, err := udppackage.UnpackUDPData(buf[:n])
		if err != nil {
			continue
		}
		go captureRemote(target, addr, src, payload, conn)
	}
}

func RunUdpServer(udpURl string) {
	// Backwards-compatible wrapper that uses Background context.
	RunUdpServerWithCtx(context.Background(), udpURl)

}

// RunUdpServerWithCtx runs the UDP server and listens for cancellation from
// the provided context. It checks the package-level Stop flag for
// compatibility with existing callers that call udpserver.Stop().
func RunUdpServerWithCtx(ctx context.Context, udpURl string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(" fatal error on udp server: ", err)
		}
	}()
	if ctx == nil {
		ctx = context.Background()
	}
	resetStopFlag()

	udpAddr, err := resolveUDPAddress(udpURl)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	serveLoop(ctx, conn)
}

func captureRemote(target, local, src *net.UDPAddr, payload []byte, conn *net.UDPConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(" fatal error on udp server: ", err)
		}
	}()

	remoteConn, err := net.DialUDP("udp", nil, target)
	if err != nil {
		log.Println("Error dialing UDP:", err)
		return
	}
	defer remoteConn.Close()

	remoteConn.SetWriteDeadline(time.Now().Add(timeout * time.Second))
	remoteConn.Write(payload)

	buf := buffer.Get()
	defer buffer.Put(buf)

	remoteConn.SetReadDeadline(time.Now().Add(timeout * time.Second))

	n, _, err := remoteConn.ReadFromUDP(buf)
	if err != nil {
		log.Println("remote failed", err)
		return
	}

	packet, err := udppackage.PackUDPData(src, target, buf[:n])
	if err != nil {
		log.Println("pack udp failed", err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(timeout * time.Second))
	conn.WriteToUDP(packet, local)
}

func Stop() {
	atomic.StoreInt32(&stop, 1)
}
