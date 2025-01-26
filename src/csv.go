package src

import (
	"encoding/csv"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"hashtechy/src/user"
	"os"
	"strconv"
	"strings"
	"sync"
)

func validateCSVInput(records [][]string) error {
	for i, record := range records {
		// Skip header row
		if i == 0 {
			continue
		}
		// Validate each field
		for _, field := range record {
			// Check for SQL injection patterns
			if strings.Contains(strings.ToLower(field), "select") ||
				strings.Contains(strings.ToLower(field), "insert") ||
				strings.Contains(strings.ToLower(field), "delete") ||
				strings.Contains(strings.ToLower(field), "update") ||
				strings.Contains(strings.ToLower(field), ";") {
				return errors.New(errors.ErrValidation, "potential SQL injection detected", nil)
			}

			// Check for XSS patterns
			if strings.Contains(field, "<script>") ||
				strings.Contains(field, "</script>") {
				return errors.New(errors.ErrValidation, "potential XSS attack detected", nil)
			}
		}
	}
	return nil
}

func parseAge(ageStr string) int {
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		logger.Error("Failed to parse age: %v", err)
		return 0
	}
	return age
}

func readCsv(name string) (err error, header []string, c chan user.User) {
	file, err := os.Open(name)
	if err != nil {
		logger.Error("Failed to open CSV file: %v", err)
		return err, nil, nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		logger.Error("Failed to read CSV file: %v", err)
		return err, nil, nil
	}

	// Validate CSV input before processing
	if err := validateCSVInput(records); err != nil {
		logger.Error("CSV validation failed: %v", err)
		return err, nil, nil
	}

	recordLength := len(records)
	if recordLength == 1 {
		err := errors.New(errors.ErrValidation, "CSV file contains only a header with no data", nil)
		logger.Error(err.Error())
		return err, nil, nil
	}

	header = records[0]
	dataRecords := records[1:]
	logger.Info("Processing CSV with %d records", len(dataRecords))

	c = make(chan user.User, len(dataRecords))
	go func() {
		defer close(c)
		var wg sync.WaitGroup

		for _, record := range dataRecords {
			wg.Add(1)
			go func(record []string) {
				defer wg.Done()

				if len(record) < 3 {
					logger.Error("Skipping record due to insufficient columns: %v", record)
					return
				}

				user := user.User{
					ID:    "",
					Email: record[0],
					Name:  record[1],
					Age:   uint8(parseAge(record[2])),
				}

				err = user.Validate()
				if err != nil {
					logger.Error("Skipping record due to validation error: %v", err)
					return
				}
				c <- user
				logger.Debug("Successfully processed user record: %s", user.Email)
			}(record)
		}

		wg.Wait()
	}()

	logger.Info("CSV processing started")
	return nil, header, c
}

// clean-up unnecessary
// func generateNewUUID() (uuid.UUID, error) {
// 	for {
// 		newID, err := uuid.NewRandom() // Generate a new UUID
// 		if err != nil {
// 			return uuid.Nil, fmt.Errorf("failed to generate UUID: %w", err) // Return uuid.Nil on error
// 		}

// 		exists, err := postgres.CheckUUIDUsed(newID)
// 		if err != nil {
// 			return uuid.Nil, fmt.Errorf("failed to check UUID usage: %w", err) // Return uuid.Nil on error
// 		}

// 		if !exists {
// 			return newID, nil // Return the unique UUID
// 		}
// 	}
// }
