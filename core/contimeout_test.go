package core

import (
	"net"
	"testing"
	"time"
)

type deadlineConn struct {
	deadline time.Time
}

func (c *deadlineConn) Read([]byte) (int, error)  { return 0, nil }
func (c *deadlineConn) Write([]byte) (int, error) { return 0, nil }
func (c *deadlineConn) Close() error              { return nil }
func (c *deadlineConn) LocalAddr() net.Addr       { return nil }
func (c *deadlineConn) RemoteAddr() net.Addr      { return nil }

func (c *deadlineConn) SetDeadline(deadline time.Time) error {
	c.deadline = deadline
	return nil
}

func (c *deadlineConn) SetReadDeadline(time.Time) error  { return nil }
func (c *deadlineConn) SetWriteDeadline(time.Time) error { return nil }

func TestSetConnectTimeoutSetsAndClearsDeadline(t *testing.T) {
	conn := &deadlineConn{}

	if err := SetConnectTimeout(conn, 1); err != nil {
		t.Fatalf("SetConnectTimeout() error = %v", err)
	}
	if conn.deadline.Before(time.Now()) {
		t.Fatal("connection deadline should be in the future")
	}

	if err := SetConnectTimeout(conn, 0); err != nil {
		t.Fatalf("SetConnectTimeout() clear error = %v", err)
	}
	if !conn.deadline.IsZero() {
		t.Fatalf("cleared deadline = %v, want zero", conn.deadline)
	}
}
