package src

import (
	"encoding/csv"
	"fmt"
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
				return fmt.Errorf("potential SQL injection detected in record %d", i)
			}

			// Check for XSS patterns
			if strings.Contains(field, "<script>") ||
				strings.Contains(field, "</script>") {
				return fmt.Errorf("potential XSS attack detected in record %d", i)
			}
		}
	}
	return nil
}

func parseAge(ageStr string) int {
	age, err := strconv.Atoi(ageStr) // Convert string to int
	if err != nil {
		return 0 // Return 0 or handle the error as needed
	}
	return age
}

func readCsv(name string) (err error, header []string, c chan user.User) {
	file, err := os.Open(name)
	if err != nil {
		return err, nil, nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err, nil, nil
	}

	// Validate CSV input before processing
	if err := validateCSVInput(records); err != nil {
		return err, nil, nil
	}

	recordLength := len(records)
	if recordLength == 1 {
		return fmt.Errorf("CSV file contains only a header with no data"), nil, nil
	}

	header = records[0]
	dataRecords := records[1:]

	c = make(chan user.User, len(dataRecords))
	go func() {
		defer close(c)
		var wg sync.WaitGroup

		for _, record := range dataRecords {
			wg.Add(1)
			go func(record []string) {
				defer wg.Done()

				if len(record) < 3 {
					fmt.Printf("Skipping record due to insufficient columns: %v\n", record)
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
					fmt.Printf("Skipping record due to validation error: %v\n", err)
					return
				}
				c <- user
			}(record)
		}

		wg.Wait()
	}()

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
