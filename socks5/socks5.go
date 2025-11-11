package socks5

import (
	"chimney3-go/core"
	"chimney3-go/settings"
	"chimney3-go/udpserver"
	"fmt"
	"log"
)

// RunServer 启动 SOCKS5 服务器或客户端，依据 isServer 参数。
func RunServer(s *settings.Settings, isServer bool) error {
	if isServer {
		return startSocks5Server(s)
	}
	return startSocks5Client(s)
}

// startSocks5Server 构建并启动 SOCKS5 服务器。
func startSocks5Server(s *settings.Settings) error {
	ss := &Socks5ServerSettings{
		ListenAddress: s.Listen,
		User:          s.Username,
		PassWord:      s.Password,
		ProxyAddress:  s.RemoteListen,
		Method:        s.Method,
	}

	go udpserver.RunUdpServer(s.Udplisten)
	log.Println("SOCKS5 server starting...")
	server := NewSocks5Server(ss, nil)
	return server.Serve()
}

// startSocks5Client 构建并启动 SOCKS5 客户端。
func startSocks5Client(s *settings.Settings) error {
	ss := &Socks5ServerSettings{
		ListenAddress: s.Listen,
		User:          s.Username,
		PassWord:      s.Password,
		ProxyAddress:  s.RemoteListen,
		Method:        s.Method,
	}
	log.Println("SOCKS5 client starting...")
	server := NewSocks5Server(ss, nil)

	httpAddr := s.Httpurl
	socks5Url := fmt.Sprintf("socks5://%s", s.Listen)
	log.Printf("Starting HTTP to SOCKS5 proxy on %s forwarding to %s", httpAddr, socks5Url)
	go core.Run2HTTP(httpAddr, socks5Url)
	return server.Serve()
}
