package all

import (
	"context"
	"fmt"

	"chimney3-go/kcpproxy"
	"chimney3-go/proxy"
	servercontext "chimney3-go/sesrvercontext"
	"chimney3-go/settings"
	"chimney3-go/socks5"
	"chimney3-go/utils"
)

type ServiceRunner func(*settings.Settings, servercontext.ServerContext) error

func serviceFactoryFor(name string) (ServiceRunner, error) {
	switch utils.NormalizeServiceName(name) {
	case SOCKS5:
		return socks5.RunServer, nil
	case PROXY:
		return proxy.RunServer, nil
	case KCP:
		return kcpproxy.RunKCPRoutine, nil
	default:
		return nil, fmt.Errorf("unknown service type: %q", name)
	}
}

const (
	PROXY  = "proxy"
	SOCKS5 = "socks5"
	KCP    = "kcp"
	SERVER = "server"
)

func ReactorWithContext(ctx context.Context, s *settings.Settings) error {
	factory, err := serviceFactoryFor(s.Which)
	if err != nil {
		return err
	}

	return factory(s, servercontext.NewServerContext(ctx))
}
