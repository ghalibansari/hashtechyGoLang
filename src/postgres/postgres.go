package postgres

import (
	"database/sql"
	"fmt"
	"hashtechy/src/user"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// Connect initializes the database connection
func Connect() error {
	connStr := "host=postgres user=postgres password=postgres dbname=crypto_db sslmode=verify-ca sslcert=/app/certs/postgres.crt sslkey=/app/certs/postgres.key sslrootcert=/app/certs/ca.crt"

	var db *sql.DB
	var err error

	// Retry logic for initial connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err = db.Ping(); err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	DB = db
	DB.SetMaxOpenConns(100)
	return nil
}

func ShowDatabases() {
	rows, err := DB.Query("SELECT datname FROM pg_database")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			panic(err)
		}
		databases = append(databases, dbName)
	}
	// Check for errors after iterating
	if err = rows.Err(); err != nil {
		panic(err)
	}
	fmt.Println("Databases:", databases)
}

func ShowTables() {
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		tables = append(tables, tableName)
	}
	// Check for errors after iterating
	if err = rows.Err(); err != nil {
		panic(err)
	}
	fmt.Println("Tables:", tables)
}

func GetAllUsers() ([]user.User, error) {
	rows, err := DB.Query("SELECT id, email, name, age FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		err := u.DecryptEmail()
		if err != nil {
			fmt.Printf("Error decrypting email for user %s: %v\n", u.ID, err)
			continue
		}
		users = append(users, u)
	}

	// Check for errors after iterating
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return users, nil
}

func GetAllUsersByQuery(name string, minAge int, maxAge int, limitParam int, skipParam int) ([]user.User, error) {
	sqlQuery := "SELECT id, email, name, age FROM users WHERE 1=1"
	var args []interface{}
	argCount := 1

	// Name filter
	if name != "" {
		name = strings.TrimSpace(name)
		if strings.ContainsAny(name, "%;'\"") {
			return nil, fmt.Errorf("invalid characters in name parameter")
		}
		sqlQuery += fmt.Sprintf(" AND name LIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	if minAge < 0 || maxAge < 0 {
		return nil, fmt.Errorf("age cannot be negative")
	}

	// Age filter (exact match or range)
	if minAge != 0 {
		sqlQuery += fmt.Sprintf(" AND age >= $%d", argCount)
		args = append(args, minAge)
		argCount++
	}

	// TODO: later move this to validate.go and use validator package
	if minAge > 120 || maxAge > 120 {
		return nil, fmt.Errorf("age exceeds maximum allowed value")
	}

	if maxAge != 0 {
		sqlQuery += fmt.Sprintf(" AND age <= $%d", argCount)
		args = append(args, maxAge)
		argCount++
	}

	// Pagination
	limit := 10 // Default limit
	if limitParam != 0 {
		limit = limitParam
	}
	sqlQuery += fmt.Sprintf(" LIMIT %d", limit)

	skip := 0
	if skipParam != 0 {
		skip = skipParam
	}
	sqlQuery += fmt.Sprintf(" OFFSET %d", skip)

	// Execute query
	rows, err := DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		err := u.DecryptEmail()
		if err != nil {
			fmt.Printf("Error decrypting email for user %s: %v\n", u.ID, err)
			continue
		}
		users = append(users, u)
	}

	// Check for errors after iterating
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return users, nil
}

func InsertUser(user user.User) (user.User, error) {
	// Encrypt email before storing
	encryptedEmail, err := user.EncryptEmail()
	if err != nil {
		return user, fmt.Errorf("failed to encrypt email: %w", err)
	}

	query := `INSERT INTO users (name, age, email) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err = DB.QueryRow(query, user.Name, user.Age, encryptedEmail).Scan(&id)
	if err != nil {
		return user, fmt.Errorf("failed to insert user: %w", err)
	}
	user.ID = id
	return user, nil
}

// clean-up unnecessary
// func CheckUUIDUsed(id uuid.UUID) (bool, error) {
// 	var exists bool
// 	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
// 	err := DB.QueryRow(query, id.String()).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check if UUID is used: %w", err)
// 	}
// 	return exists, nil
// }

// GenerateUUID retrieves a new UUID from the database // clean-up unnecessary
// func GenerateUUID() (string, error) {
// 	var newID string
// 	err := DB.QueryRow("SELECT uuid_generate_v4()").Scan(&newID)
// 	if err != nil || newID == "" {
// 		return "", fmt.Errorf("failed to generate a valid UUID: %w", err)
// 	}
// 	return newID, nil
// }
