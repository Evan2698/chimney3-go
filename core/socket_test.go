package core

import (
	"bytes"
	"chimney3-go/privacy"
	"net"
	"testing"
	"time"
)

type mockCon struct {
	k bytes.Buffer
}

func (m *mockCon) Read(b []byte) (n int, err error) {

	return m.k.Read(b)
}

func (m *mockCon) Write(b []byte) (n int, err error) {

	return m.k.Write(b)
}

func (m *mockCon) Close() error {

	return nil
}

func (m *mockCon) LocalAddr() net.Addr {

	return nil
}

func (m *mockCon) RemoteAddr() net.Addr {

	return nil
}

// A zero value for t means I/O operations will not time out.
func (m *mockCon) SetDeadline(t time.Time) error {

	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (m *mockCon) SetReadDeadline(t time.Time) error {

	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (m *mockCon) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestSocket(t *testing.T) {
	con := &mockCon{}

	socket := NewSocks5Socket(con, privacy.NewMethodWithName("CHACHA-20"),
		privacy.MakeCompressKey("hello"), nil, nil)

	socket.Write([]byte("hello world!!!"))

	k := make([]byte, 90)
	n, err := socket.Read(k)
	t.Log(n, k[:n], err)
	t.Log(string(k[:n]))

}

func TestSSLListenerCloseIsIdempotent(t *testing.T) {
	l := &SSLListenerImpl{ListenChannel: make(chan MySSLSocket), RawListener: &net.TCPListener{}}
	if err := l.Close(); err == nil {
		t.Fatal("Close() should return an error for a nil or unusable underlying listener")
	}
	if err := l.Close(); err == nil {
		t.Fatal("Close() should be safe to call twice")
	}
}

func TestSocks5AddressNilReceiverIsSafe(t *testing.T) {
	var addr *Socks5Address
	if got := addr.GetAddress(); got != "" {
		t.Fatalf("nil receiver GetAddress() = %q, want empty string", got)
	}
	if got := addr.String(); got != "" {
		t.Fatalf("nil receiver String() = %q, want empty string", got)
	}
	addr.SetIPv4Address([]byte{127, 0, 0, 1}, 8080)
	if addr != nil {
		t.Fatal("nil receiver should not be mutated")
	}
}
