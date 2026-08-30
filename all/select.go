package all

import (
	"context"
	"fmt"
	"strings"

	"chimney3-go/kcpproxy"
	"chimney3-go/proxy"
	"chimney3-go/settings"
	"chimney3-go/socks5"
)

type ServiceRunner func(*settings.Settings, bool) error

func normalizeServiceName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		return SOCKS5
	}
	return normalized
}

func normalizeRuntimeMode(mode string) string {
	return strings.ToLower(strings.TrimSpace(mode))
}

func isServerMode(mode string) bool {
	normalized := normalizeRuntimeMode(mode)
	return normalized != "client"
}

func serviceFactoryFor(name string) (ServiceRunner, error) {
	switch normalizeServiceName(name) {
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

// Reactor selects and starts the appropriate subsystem based on the
// provided configuration. It returns an error when the selection is
// unknown or when the selected subsystem reports an error.
func Reactor(s *settings.Settings) error {
	return ReactorWithContext(context.Background(), s)
}

func ReactorWithContext(ctx context.Context, s *settings.Settings) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil {
		return fmt.Errorf("settings: nil")
	}

	factory, err := serviceFactoryFor(s.Which)
	if err != nil {
		return err
	}

	isServer := isServerMode(s.Mode)
	return factory(s, isServer)
}
