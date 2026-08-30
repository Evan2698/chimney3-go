package socks5

import (
	"io"
	"net"
	"testing"
)

func TestSocks5CloseClosesUnderlyingConn(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	s := &Socks5{conn: clientConn}
	s.Close()

	buf := make([]byte, 1)
	_, err := serverConn.Read(buf)
	if err == nil || err != io.EOF {
		t.Fatalf("expected EOF after close, got %v", err)
	}
}

func TestRunServerRejectsNilSettings(t *testing.T) {
	if err := RunServer(nil, true); err == nil {
		t.Fatal("RunServer(nil, true) should reject nil settings")
	}
	if err := RunServer(nil, false); err == nil {
		t.Fatal("RunServer(nil, false) should reject nil settings")
	}

	server := &Socks5S{}
	if err := server.Serve(); err == nil {
		t.Fatal("Socks5S.Serve() should reject nil settings")
	}
}

func TestSocks5DialRejectsNilReceiver(t *testing.T) {
	var client *Socks5
	if _, err := client.Dial(nil); err == nil {
		t.Fatal("Socks5.Dial on nil receiver should reject nil")
	}
}

func TestValidateSocks5Method(t *testing.T) {
	if err := validateSocks5Method("unknown"); err == nil {
		t.Fatal("validateSocks5Method should reject unsupported method")
	}
	if err := validateSocks5Method("CHACHA-20"); err != nil {
		t.Fatalf("validateSocks5Method should accept registered method: %v", err)
	}
}
