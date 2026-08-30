package utils

import (
	"chimney3-go/udpserver"
)

func StartUDPServerIfConfigured(addr string) {
	if addr == "" {
		return
	}
	go udpserver.RunUdpServer(addr)
}
