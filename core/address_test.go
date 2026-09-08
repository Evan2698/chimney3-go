package core

import "testing"

func TestSocks5AddressParseRejectsMalformedInput(t *testing.T) {
	for _, input := range [][]byte{
		{ADDRESSTYPE_DOMAIN, 10, 'x'},
		{ADDRESSTYPE_IPV4, 127, 0, 0},
		{ADDRESSTYPE_IPV6, 0, 0, 0},
	} {
		address := NewSocks5Address()
		if err := address.Parse(input); err == nil {
			t.Fatalf("Parse(%v) succeeded for malformed input", input)
		}
		if address.Valid {
			t.Fatalf("Parse(%v) marked malformed input valid", input)
		}
	}
}
