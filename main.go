package main

import (
	"hashtechy/src"
	"hashtechy/src/logger"
	"os"
)

func main() {
	err := src.App()
	if err != nil {
		logger.Error("Application error: %v", err)
		os.Exit(1)
	}
}
