package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"chimney3-go/all"
	"chimney3-go/settings"
	"chimney3-go/utils"
)

func resolveConfigPath(dir string) string {
	return filepath.Join(dir, "configs", "setting.json")
}

type App struct {
	configPath string
	cfg        *settings.Settings
}

func NewApp(configPath string) (*App, error) {
	if strings.TrimSpace(configPath) == "" {
		return nil, errors.New("config path is empty")
	}
	cfg, err := settings.Parse(configPath)
	if err != nil {
		return nil, err
	}
	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &App{configPath: configPath, cfg: cfg}, nil
}

func (a *App) Run(ctx context.Context) error {
	if a == nil {
		return errors.New("app: nil")
	}
	if a.cfg == nil {
		return errors.New("app config: nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return all.ReactorWithContext(ctx, a.cfg)
}

func waitForShutdown(ctx context.Context, done <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, ok := <-done:
		if !ok {
			return nil
		}
		return err
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("crash on :", r)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	dir, err := utils.RetrieveExePath()
	if err != nil {
		log.Fatalf("failed to determine executable path: %v", err)
	}

	cfgFlag := flag.String("config", "", "path to JSON config file (default: <exe>/configs/setting.json)")
	flag.Parse()

	jsonPath := *cfgFlag
	if jsonPath == "" {
		jsonPath = resolveConfigPath(dir)
	}

	app, err := NewApp(jsonPath)
	if err != nil {
		log.Fatalf("failed to load config %s: %v", jsonPath, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("starting chimney3-go; config=%s", jsonPath)

	done := make(chan error, 1)
	go func() {
		done <- app.Run(ctx)
	}()

	err = waitForShutdown(ctx, done)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("shutdown signal received; exiting")
			os.Exit(0)
		}
		log.Fatalf("service exited with error: %v", err)
	}
}
