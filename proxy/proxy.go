package proxy

import (
	"chimney3-go/settings"
	"chimney3-go/utils"
	"fmt"
)

// RunServer starts the proxy subsystem. It returns an error when startup or
// the running subsystem encounters a terminal error. Callers should decide
// whether to log/fatal or attempt recovery.
func RunServer(s *settings.Settings, isServer bool) error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	if isServer {
		return runserver(s)
	}
	return runclient(s)
}

func runclient(s *settings.Settings) error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	pc := &proxyClient{
		Password:     s.Password,
		LocalHost:    s.Listen,
		ProxyAddress: s.RemoteListen,
		Exit:         false,
	}
	return pc.Serve()
}

func runserver(s *settings.Settings) error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	ps := &proxyServer{
		Host:     s.Listen,
		Password: s.Password,
		Which:    s.Which,
		Exit:     false,
	}

	utils.StartUDPServerIfConfigured(s.Udplisten)
	return ps.Serve()
}
