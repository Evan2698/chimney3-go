package main

import (
	"chimney3-go/all"
	"chimney3-go/settings"
	"chimney3-go/utils"
	"fmt"
	"os"
	"runtime"
)

var (
	isServer *bool
)

func main() {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu * 4)

	dir, err := utils.RetrieveExePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to determine executable path:", err)
		os.Exit(1)
	}

	jsonPath := dir + "configs/setting.json"
	settings, err := settings.Parse(jsonPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config:", err)
		os.Exit(1)
	}
	all.Reactor(settings)
}
