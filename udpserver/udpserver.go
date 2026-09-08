package udpserver

import (
	servercontext "chimney3-go/sesrvercontext"
	"log"
	"net"
	"sync"
	"time"
	"tun2proxylib/gvisorcore/buffer"
	"tun2proxylib/udppackage"
)

const (
	defaultUDPURL = "0.0.0.0:5353"
	timeout       = 20 // seconds
)

func resolveUDPAddress(udpURL string) (*net.UDPAddr, error) {
	if udpURL == "" {
		udpURL = defaultUDPURL
	}
	return net.ResolveUDPAddr("udp", udpURL)
}

func serveLoop(ctx servercontext.ServerContext, conn *net.UDPConn, writeMu *sync.Mutex) {
	buf := buffer.Get()
	defer buffer.Put(buf)

	for {
		if ctx.IsInterrupted() {
			return
		}

		if err := conn.SetReadDeadline(time.Now().Add(timeout * time.Second)); err != nil {
			log.Println("set UDP read deadline failed:", err)
			continue
		}
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
				log.Println("read UDP packet failed:", err)
				continue
			}
			continue
		}

		target, src, payload, err := udppackage.UnpackUDPData(buf[:n])
		if err != nil {
			log.Println("unpack UDP data failed:", err)
			continue
		}
		go captureRemote(target, addr, src, payload, conn, writeMu)
	}
}

func RunUdpServer(udpURl string, ctx servercontext.ServerContext) {
	// Backwards-compatible wrapper that uses Background context.
	RunUdpServerWithCtx(ctx, udpURl)

}

// RunUdpServerWithCtx runs the UDP server and listens for cancellation from
// the provided context. It checks the package-level Stop flag for
// compatibility with existing callers that call udpserver.Stop().
func RunUdpServerWithCtx(ctx servercontext.ServerContext, udpURl string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(" fatal error on udp server: ", err)
		}
	}()

	ctx.ClearInterrupted()

	udpAddr, err := resolveUDPAddress(udpURl)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	var writeMu sync.Mutex
	serveLoop(ctx, conn, &writeMu)
}

func captureRemote(target, local, src *net.UDPAddr, payload []byte, conn *net.UDPConn, writeMu *sync.Mutex) {
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

	if err := remoteConn.SetWriteDeadline(time.Now().Add(timeout * time.Second)); err != nil {
		log.Println("set remote UDP write deadline failed:", err)
		return
	}
	if _, err := remoteConn.Write(payload); err != nil {
		log.Println("write remote UDP packet failed:", err)
		return
	}

	buf := buffer.Get()
	defer buffer.Put(buf)

	if err := remoteConn.SetReadDeadline(time.Now().Add(timeout * time.Second)); err != nil {
		log.Println("set remote UDP read deadline failed:", err)
		return
	}

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

	writeMu.Lock()
	defer writeMu.Unlock()
	if err := conn.SetWriteDeadline(time.Now().Add(timeout * time.Second)); err != nil {
		log.Println("set UDP write deadline failed:", err)
		return
	}
	if _, err := conn.WriteToUDP(packet, local); err != nil {
		log.Println("write UDP response failed:", err)
	}
}
