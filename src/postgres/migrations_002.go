package postgres

import "hashtechy/src/logger"

func CreateIndexes() error {
	// Create index on name for faster search
	_, err := DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
	`)
	if err != nil {
		logger.Error("Failed to create name index: %v", err)
		return err
	}

	// Create index on age for faster filtering
	_, err = DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_age ON users(age);
	`)
	if err != nil {
		logger.Error("Failed to create age index: %v", err)
		return err
	}

	logger.Info("Successfully created database indexes")
	return nil
}
