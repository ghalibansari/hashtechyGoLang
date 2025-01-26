package tests

import (
	"hashtechy/src/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	err := postgres.Connect()
	assert.NoError(t, err, "Database connection should succeed")
	defer postgres.DB.Close()
}

func TestCreateUserTable(t *testing.T) {
	err := postgres.Connect()
	assert.NoError(t, err, "Database connection should succeed")
	defer postgres.DB.Close()

	postgres.DropUserTable()
	err = postgres.CreateUserTable()
	assert.NoError(t, err, "User table creation should succeed")
}

func TestCreateIndexes(t *testing.T) {
	err := postgres.Connect()
	assert.NoError(t, err, "Database connection should succeed")
	defer postgres.DB.Close()

	err = postgres.CreateIndexes()
	assert.NoError(t, err, "Index creation should succeed")
}
