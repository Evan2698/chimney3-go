package proxy

import (
	"chimney3-go/settings"
	"chimney3-go/udpserver"
)

// RunServer starts the proxy subsystem. It returns an error when startup or
// the running subsystem encounters a terminal error. Callers should decide
// whether to log/fatal or attempt recovery.
func RunServer(s *settings.Settings, isServer bool) error {
	if isServer {
		return runserver(s)
	}
	return runclient(s)
}

func runclient(s *settings.Settings) error {
	pc := &proxyClient{
		Password:     s.Password,
		LocalHost:    s.Listen,
		ProxyAddress: s.RemoteListen,
		Exit:         false,
	}
	return pc.Serve()
}

func runserver(s *settings.Settings) error {
	ps := &proxyServer{
		Host:     s.Listen,
		Password: s.Password,
		Which:    s.Which,
		Exit:     false,
	}

	udpServerAddr := s.Udplisten
	go udpserver.RunUdpServer(udpServerAddr)
	return ps.Serve()
}
