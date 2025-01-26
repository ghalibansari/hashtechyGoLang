package main

import (
	"fmt"
	"hashtechy/src"
	"os"
)

func main() {
	err := src.App()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
