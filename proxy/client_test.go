package proxy

import (
	"testing"
	"time"
)

func TestProxyClientCloseStopsServe(t *testing.T) {
	c := &proxyClient{
		Password:     "test-password",
		LocalHost:    "127.0.0.1:0",
		ProxyAddress: "127.0.0.1:1",
	}

	serveDone := make(chan error, 1)
	go func() {
		serveDone <- c.Serve()
	}()

	time.Sleep(50 * time.Millisecond)

	if err := c.Close(); err != nil {
		t.Fatalf("Close() returned unexpected error: %v", err)
	}

	select {
	case err := <-serveDone:
		if err != nil {
			t.Fatalf("Serve() returned unexpected error after Close(): %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Serve() did not exit after Close()")
	}
}

func TestRunServerRejectsNilSettings(t *testing.T) {
	if err := RunServer(nil, true); err == nil {
		t.Fatal("RunServer(nil, true) should reject nil settings")
	}
	if err := RunServer(nil, false); err == nil {
		t.Fatal("RunServer(nil, false) should reject nil settings")
	}
}

func TestProxyServerServeRejectsNilReceiver(t *testing.T) {
	var server *proxyServer
	if err := server.Serve(); err == nil {
		t.Fatal("proxyServer.Serve on nil receiver should return an error")
	}
}
