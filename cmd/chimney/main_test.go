package main

import (
	"context"
	"errors"
	"testing"
)

func TestWaitForShutdownReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitForShutdown(ctx, make(chan error))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitForShutdown() = %v, want %v", err, context.Canceled)
	}
}

func TestWaitForShutdownReturnsServiceError(t *testing.T) {
	done := make(chan error, 1)
	expected := errors.New("boom")
	done <- expected

	err := waitForShutdown(context.Background(), done)
	if !errors.Is(err, expected) {
		t.Fatalf("waitForShutdown() = %v, want %v", err, expected)
	}
}

func TestAppRunRejectsNilConfig(t *testing.T) {
	var app *App
	if err := app.Run(context.Background()); err == nil {
		t.Fatal("App.Run on nil receiver should return an error")
	}

	app = &App{}
	if err := app.Run(context.Background()); err == nil {
		t.Fatal("App.Run with nil cfg should return an error")
	}
}
