package postgres

import (
	"hashtechy/src/errors"
	"hashtechy/src/logger"
)

func CreateUserTable() error {
	// Enable the uuid-ossp extension
	_, err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		logger.Error("Failed to create extension: %v", err)
		return errors.New(errors.ErrDatabase, "failed to create uuid-ossp extension", err)
	}
	logger.Info("Successfully created or verified uuid-ossp extension")

	// Create the user table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name TEXT,
			age INT,
			email TEXT UNIQUE
		);
	`)
	if err != nil {
		logger.Error("Failed to create user table: %v", err)
		return errors.New(errors.ErrDatabase, "failed to create user table", err)
	}
	logger.Info("Successfully created or verified users table")
	return nil
}

// Caution: Ensure that the user table is not in use before dropping it to avoid data loss.
func DropUserTable() error {
	_, err := DB.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		logger.Error("Failed to drop user table: %v", err)
		return errors.New(errors.ErrDatabase, "failed to drop user table", err)
	}
	logger.Info("Successfully dropped users table")
	return nil
}
