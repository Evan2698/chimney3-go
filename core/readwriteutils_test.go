package core

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestReadXBytesRejectsShortBuffer(t *testing.T) {
	if _, err := ReadXBytes(2, make([]byte, 1), nil); err == nil {
		t.Fatal("ReadXBytes accepted a short buffer")
	}
}

type noProgressConn struct{}

func (noProgressConn) Read([]byte) (int, error)           { return 0, nil }
func (noProgressConn) Write([]byte) (int, error)          { return 0, errors.New("unused") }
func (noProgressConn) Close() error                       { return nil }
func (noProgressConn) LocalAddr() net.Addr                { return nil }
func (noProgressConn) RemoteAddr() net.Addr               { return nil }
func (noProgressConn) SetDeadline(_ time.Time) error      { return nil }
func (noProgressConn) SetReadDeadline(_ time.Time) error  { return nil }
func (noProgressConn) SetWriteDeadline(_ time.Time) error { return nil }

func TestReadXBytesRejectsNoProgress(t *testing.T) {
	if _, err := ReadXBytes(1, make([]byte, 1), noProgressConn{}); !errors.Is(err, io.ErrNoProgress) {
		t.Fatalf("ReadXBytes error = %v, want io.ErrNoProgress", err)
	}
}
