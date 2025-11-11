package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"chimney3-go/all"
	"chimney3-go/settings"
	"chimney3-go/utils"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	dir, err := utils.RetrieveExePath()
	if err != nil {
		log.Fatalf("failed to determine executable path: %v", err)
	}

	// Allow user to pass a config path; default to executable dir + /configs/setting.json
	cfgFlag := flag.String("config", "", "path to JSON config file (default: <exe>/configs/setting.json)")
	flag.Parse()

	jsonPath := *cfgFlag
	if jsonPath == "" {
		jsonPath = dir + "/configs/setting.json"
	}

	cfg, err := settings.Parse(jsonPath)
	if err != nil {
		log.Fatalf("failed to load config %s: %v", jsonPath, err)
	}

	// validate configuration early
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Setup signal handling for a graceful shutdown notice.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// startup logging
	log.Printf("starting chimney3-go; config=%s", jsonPath)

	// Run the selected reactor. It's blocking; when it returns we exit.
	done := make(chan error, 1)
	go func() {
		done <- all.Reactor(cfg)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received; exiting")
		os.Exit(0)
	case err := <-done:
		if err != nil {
			log.Fatalf("service exited with error: %v", err)
		}
	}
}
