package tests

import (
	"hashtechy/src"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullApplicationFlow(t *testing.T) {
	// This test simulates the entire application flow
	err := src.App()
	assert.NoError(t, err, "Full application should run without critical errors")

	// Add more specific assertions about the application state
	// For example, check database connections, server status, etc.
}
