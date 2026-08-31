package utils

import (
	servercontext "chimney3-go/sesrvercontext"
	"chimney3-go/udpserver"
	"context"
	"fmt"
)

func StartUDPServerIfConfigured(addr string) (servercontext.ServerContext, error) {
	if addr == "" {
		return nil, fmt.Errorf("UDP address is not configured")
	}
	ctx := servercontext.NewServerContext(context.Background())
	go udpserver.RunUdpServer(addr, ctx)
	return ctx, nil
}
