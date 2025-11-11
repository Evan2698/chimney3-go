package all

import (
	"fmt"
	"strings"

	"chimney3-go/kcpproxy"
	"chimney3-go/proxy"
	"chimney3-go/settings"
	"chimney3-go/socks5"
)

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
	isServer := strings.EqualFold(s.Mode, SERVER)

	switch strings.ToLower(s.Which) {
	case SOCKS5:
		return socks5.RunServer(s, isServer)
	case PROXY:
		// proxy.RunServer currently doesn't return an error; keep
		// behavior but wrap in nil for a consistent signature.
		proxy.RunServer(s, isServer)
		return nil
	case KCP:
		return kcpproxy.RunKCPRoutine(s, isServer)
	default:
		return fmt.Errorf("unknown service type: %q", s.Which)
	}
}
