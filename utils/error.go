package utils

import (
	"fmt"
	"log"
)

func WrapError(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func LogError(action string, err error) {
	if err == nil {
		return
	}
	log.Printf("%s: %v", action, err)
}

func Recover(action string) {
	if r := recover(); r != nil {
		log.Printf("%s: recovered panic: %v", action, r)
	}
}
