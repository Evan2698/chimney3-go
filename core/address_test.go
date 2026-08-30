package core

import "testing"

func TestParseTargetAddressIPv4UsesFourByteRepresentation(t *testing.T) {
	addr, err := ParseTargetAddress("127.0.0.1:8080")
	if err != nil {
		t.Fatalf("ParseTargetAddress returned error: %v", err)
	}
	if addr.Type != ADDRESSTYPE_IPV4 {
		t.Fatalf("expected IPv4 type, got %d", addr.Type)
	}
	if len(addr.IP) != 4 {
		t.Fatalf("expected 4-byte IPv4 representation, got %d bytes", len(addr.IP))
	}

	got := addr.Bytes()
	if len(got) != 7 {
		t.Fatalf("expected 7-byte wire format, got %d bytes: %v", len(got), got)
	}
	if got[0] != ADDRESSTYPE_IPV4 {
		t.Fatalf("expected type byte 0x01, got 0x%02x", got[0])
	}
}

func TestSocks5AddressRoundTripForIPv4(t *testing.T) {
	addr := NewSocks5Address()
	addr.SetIPv4Address([]byte{127, 0, 0, 1}, 8080)

	wire := addr.Bytes()
	parsed := NewSocks5Address()
	if err := parsed.Parse(wire); err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if parsed.Type != ADDRESSTYPE_IPV4 {
		t.Fatalf("expected IPv4 type after parse, got %d", parsed.Type)
	}
	if string(parsed.IP) != string([]byte{127, 0, 0, 1}) {
		t.Fatalf("expected original IPv4 bytes, got %v", parsed.IP)
	}
}
