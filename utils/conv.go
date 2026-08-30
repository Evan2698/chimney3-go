package utils

import (
	"encoding/binary"
	"strings"
)

// Int2Bytes converts a uint32 to a 4-byte big-endian slice.
func Int2Bytes(n uint32) []byte {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, n)
	return out
}

// Bytes2Int converts a 4-byte big-endian slice to uint32.
func Bytes2Int(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return binary.BigEndian.Uint32(b[:4])
}

// Bytes2Uint16 converts a 2-byte big-endian slice to uint16.
func Bytes2Uint16(b []byte) uint16 {
	if len(b) < 2 {
		return 0
	}
	return binary.BigEndian.Uint16(b[:2])
}

// Uint162Bytes converts a uint16 to a 2-byte big-endian slice.
func Uint162Bytes(b uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, b)
	return out
}

// Port2Bytes returns the big-endian encoding of a port.
func Port2Bytes(port uint16) []byte {
	return Uint162Bytes(port)
}

// FormatProtocol normalizes the protocol name and prefers QUIC when present.
func FormatProtocol(p string) string {
	normalized := strings.ToLower(strings.TrimSpace(p))
	if strings.Contains(normalized, "quic") {
		return "quic"
	}
	return "tcp"
}
