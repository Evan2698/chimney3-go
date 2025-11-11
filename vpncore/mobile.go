package vpncore

import (
	"chimney3-go/socks5"
	"context"
	"log"

	"tun2proxylib/mobile"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type Protect interface {
	mobile.ProtectSocket
}

type Chimney struct {
	Fd          int
	Pfun        Protect
	User        string
	Pass        string
	TcpProxyUrl string
	MTU         int
	UdpProxyUrl string
}

var (
	client   socks5.Socks5Server
	netstack *stack.Stack
)

// StartChimney starts the VPN components and runs the socks5 server under
// the given context. When ctx is canceled, the server and netstack will be
// stopped.
func StartChimney(ctx context.Context, c *Chimney) error {

	var err error
	client = buildVpnClient("127.0.0.1:1080", c.TcpProxyUrl, c.User, c.Pass, c.Pfun)
	netstack, err = buildNetstackVpnClient(c.Fd, uint32(c.MTU), "127.0.0.1:1080", c.UdpProxyUrl, c.Pfun)
	if err != nil {
		return err
	}

	// run the socks5 server in the background and stop it when ctx is done
	go func() {
		if client == nil {
			return
		}
		if err := client.Serve(); err != nil {
			log.Println("socks5 server exited with error:", err)
		}
	}()

	go func() {
		<-ctx.Done()
		if client != nil {
			client.Stop()
			client = nil
		}
		if netstack != nil {
			netstack.Close()
			netstack = nil
		}
	}()

	return nil
}

func StopChimney() {

	if netstack != nil {
		netstack.Close()
		netstack = nil
	}

	if client != nil {
		client.Stop()
		client = nil
	}
}
