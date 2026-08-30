package all

import (
	"context"
	"errors"
	"testing"

	"chimney3-go/settings"
)

func TestNormalizeServiceName(t *testing.T) {
	cases := map[string]string{
		"SOCKS5": SOCKS5,
		"Proxy":  PROXY,
		"KCP":    KCP,
		"server": SERVER,
		"":       SOCKS5,
	}

	for input, want := range cases {
		if got := normalizeServiceName(input); got != want {
			t.Fatalf("normalizeServiceName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestIsServerMode(t *testing.T) {
	cases := map[string]bool{
		"":         true,
		"server":   true,
		"SERVER":   true,
		" client ": false,
		"client":   false,
	}
	for input, want := range cases {
		if got := isServerMode(input); got != want {
			t.Fatalf("isServerMode(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestServiceFactoryForUnknownName(t *testing.T) {
	if _, err := serviceFactoryFor("unknown"); err == nil {
		t.Fatal("serviceFactoryFor should reject unknown service name")
	}
}

func TestReactorWithContextReturnsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := ReactorWithContext(ctx, &settings.Settings{Which: PROXY, Mode: SERVER})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
