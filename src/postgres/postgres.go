package postgres

import (
	"database/sql"
	"fmt"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"hashtechy/src/user"
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
			logger.Error("Attempt %d: Failed to open database: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err = db.Ping(); err == nil {
			break
		}
		logger.Error("Attempt %d: Failed to connect to database: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return errors.New(errors.ErrDatabase, "failed to connect to database after retries", err)
	}

	DB = db
	DB.SetMaxOpenConns(100)
	logger.Info("Successfully connected to database")
	return nil
}

func ShowDatabases() {
	rows, err := DB.Query("SELECT datname FROM pg_database")
	if err != nil {
		logger.Error("Failed to query databases: %v", err)
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			logger.Error("Failed to scan database name: %v", err)
			return
		}
		databases = append(databases, dbName)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error during database iteration: %v", err)
		return
	}
	logger.Info("Available databases: %v", databases)
}

func ShowTables() {
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		logger.Error("Failed to query tables: %v", err)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			logger.Error("Failed to scan table name: %v", err)
			return
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error during table iteration: %v", err)
		return
	}
	logger.Info("Available tables: %v", tables)
}

func GetAllUsers() ([]user.User, error) {
	rows, err := DB.Query("SELECT id, email, name, age FROM users")
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, "failed to query users", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age); err != nil {
			return nil, errors.New(errors.ErrDatabase, "failed to scan user", err)
		}
		err := u.DecryptEmail()
		if err != nil {
			logger.Error("Error decrypting email for user %s: %v", u.ID, err)
			continue
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New(errors.ErrDatabase, "error during rows iteration", err)
	}

	logger.Debug("Retrieved %d users from database", len(users))
	return users, nil
}

func GetAllUsersByQuery(name string, minAge int, maxAge int, limitParam int, skipParam int) ([]user.User, error) {
	sqlQuery := "SELECT id, email, name, age FROM users WHERE 1=1"
	var args []interface{}
	argCount := 1

	if name != "" {
		name = strings.TrimSpace(name)
		if strings.ContainsAny(name, "%;'\"") {
			return nil, errors.New(errors.ErrValidation, "invalid characters in name parameter", nil)
		}
		sqlQuery += fmt.Sprintf(" AND name LIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	if minAge < 0 || maxAge < 0 {
		return nil, errors.New(errors.ErrValidation, "age cannot be negative", nil)
	}

	if minAge != 0 {
		sqlQuery += fmt.Sprintf(" AND age >= $%d", argCount)
		args = append(args, minAge)
		argCount++
	}

	if minAge > 120 || maxAge > 120 {
		return nil, errors.New(errors.ErrValidation, "age exceeds maximum allowed value", nil)
	}

	if maxAge != 0 {
		sqlQuery += fmt.Sprintf(" AND age <= $%d", argCount)
		args = append(args, maxAge)
		argCount++
	}

	limit := 10
	if limitParam != 0 {
		limit = limitParam
	}
	sqlQuery += fmt.Sprintf(" LIMIT %d", limit)

	skip := 0
	if skipParam != 0 {
		skip = skipParam
	}
	sqlQuery += fmt.Sprintf(" OFFSET %d", skip)

	logger.Debug("Executing query: %s with args: %v", sqlQuery, args)
	rows, err := DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, "failed to query users", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age); err != nil {
			return nil, errors.New(errors.ErrDatabase, "failed to scan user", err)
		}
		err := u.DecryptEmail()
		if err != nil {
			logger.Error("Error decrypting email for user %s: %v", u.ID, err)
			continue
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New(errors.ErrDatabase, "error during rows iteration", err)
	}

	logger.Debug("Retrieved %d users from query", len(users))
	return users, nil
}

func InsertUser(user user.User) (user.User, error) {
	encryptedEmail, err := user.EncryptEmail()
	if err != nil {
		return user, errors.New(errors.ErrEncryption, "failed to encrypt email", err)
	}

	query := `INSERT INTO users (name, age, email) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err = DB.QueryRow(query, user.Name, user.Age, encryptedEmail).Scan(&id)
	if err != nil {
		return user, errors.New(errors.ErrDatabase, "failed to insert user", err)
	}
	user.ID = id
	logger.Info("Successfully inserted user with ID: %s", id)
	return user, nil
}
