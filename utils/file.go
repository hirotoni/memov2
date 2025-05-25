package utils

import (
	"errors"
	"log"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		log.Printf("Error checking file existence: %v\n", err)
		return false
	}
	return true
}
