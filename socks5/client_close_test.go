package socks5

import (
	"io"
	"net"
	"testing"
)

func TestSocks5CloseClosesUnderlyingConn(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	s := &Socks5{Conn: clientConn}
	s.Close()

	buf := make([]byte, 1)
	_, err := serverConn.Read(buf)
	if err == nil || err != io.EOF {
		t.Fatalf("expected EOF after close, got %v", err)
	}
}

func TestRunServerRejectsNilSettings(t *testing.T) {

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
