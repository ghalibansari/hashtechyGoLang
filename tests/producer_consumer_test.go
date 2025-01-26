package tests

import (
	"testing"

	"hashtechy/src"

	"github.com/stretchr/testify/assert"
)

func TestProducerConsumer(t *testing.T) {
	// Mock or use a test-specific configuration
	err := src.Producer()
	assert.NoError(t, err, "Producer should run without errors")

	// You might need to add more specific tests based on your implementation
	// Consider mocking RabbitMQ for controlled testing
}
