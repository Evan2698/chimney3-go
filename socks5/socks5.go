package socks5

import (
	"chimney3-go/core"
	servercontext "chimney3-go/sesrvercontext"
	"chimney3-go/settings"
	"chimney3-go/utils"
	"fmt"
	"log"
)

func buildSocks5Settings(cfg *settings.Settings) *Socks5ServerSettings {
	if cfg == nil {
		return &Socks5ServerSettings{}
	}
	return &Socks5ServerSettings{
		ListenAddress: cfg.Listen,
		User:          cfg.Username,
		PassWord:      cfg.Password,
		ProxyAddress:  cfg.RemoteListen,
		Method:        cfg.Method,
	}
}

func startHTTPBridge(httpAddr, socks5Addr string) {
	if httpAddr == "" || socks5Addr == "" {
		return
	}
	log.Printf("Starting HTTP to SOCKS5 proxy on %s forwarding to %s", httpAddr, socks5Addr)
	go core.Run2HTTP(httpAddr, socks5Addr)
}

// RunServer 启动 SOCKS5 服务器或客户端，依据 isServer 参数。
func RunServer(s *settings.Settings, ctx servercontext.ServerContext) error {
	isServer := utils.IsServerMode(s.Mode)
	if isServer {
		return startSocks5Server(s, ctx)
	}
	return startSocks5Client(s, ctx)
}

// startSocks5Server 构建并启动 SOCKS5 服务器。
func startSocks5Server(s *settings.Settings, ctx servercontext.ServerContext) error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	ss := buildSocks5Settings(s)
	utils.StartUDPServerIfConfigured(s.Udplisten)
	log.Println("SOCKS5 server starting...")
	server := NewSocks5Server(ss, nil, ctx)
	return server.Serve()
}

// startSocks5Client 构建并启动 SOCKS5 客户端。
func startSocks5Client(s *settings.Settings, ctx servercontext.ServerContext) error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	ss := buildSocks5Settings(s)
	log.Println("SOCKS5 client starting...")
	server := NewSocks5Server(ss, nil, ctx)

	httpAddr := s.Httpurl
	socks5URL := fmt.Sprintf("socks5://%s", s.Listen)
	startHTTPBridge(httpAddr, socks5URL)
	return server.Serve()
}
