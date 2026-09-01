package all

import (
	"context"
	"errors"
	"testing"

	"chimney3-go/settings"
)

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
