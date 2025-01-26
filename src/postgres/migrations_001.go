package postgres

import "log"

func CreateUserTable() {
	// Enable the uuid-ossp extension
	_, err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		log.Fatalf("Failed to create extension: %v", err)
	}

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
		log.Fatalf("Failed to create user table: %v", err)
	}
}

// Caution: Ensure that the user table is not in use before dropping it to avoid data loss.
func DropUserTable() {
	_, err := DB.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		panic(err)
	}
}
