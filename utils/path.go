package utils

import (
	"log"
	"os"
	"path/filepath"
)

func RetrieveExePath() (string, error) {
	file, err := os.Executable()
	if err != nil {
		return "", err
	}
	re, err := filepath.Abs(file)
	if err != nil {
		log.Printf("failed to get absolute path of executable: %v", err)
	}
	//log.Print("The path is ", re)
	return filepath.Dir(re), err
}

func RetrieveCertsPath() (string, error) {
	exePath, err := RetrieveExePath()
	if err != nil {
		return "", err
	}

	exePath = exePath + "/certs"
	return exePath, err
}
