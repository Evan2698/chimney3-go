package proxy

import (
	"chimney3-go/settings"
	"chimney3-go/udpserver"
)

func RunServer(s *settings.Settings, isServer bool) {

	if isServer {
		runserver(s)

	} else {
		runclient(s)
	}
}

func runclient(s *settings.Settings) {
	pc := &proxyClient{
		Password:     s.Password,
		LocalHost:    s.Listen,
		ProxyAddress: s.RemoteListen,
		Exit:         false,
	}
	pc.Serve()
}

func runserver(s *settings.Settings) {
	ps := &proxyServer{
		Host:     s.Listen,
		Password: s.Password,
		Which:    s.Which,
		Exit:     false,
	}

	udpServerAddr := s.Udplisten
	go udpserver.RunUdpServer(udpServerAddr)
	ps.Serve()
}
