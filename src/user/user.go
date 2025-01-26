package user

import (
	"fmt"
	"regexp"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   uint8  `json:"age"`
}

func (u *User) Validate() error {
	if len(u.Name) < 2 || len(u.Name) > 50 {
		return fmt.Errorf("name must be between 2 and 50 characters")
	}

	if u.Age < 0 || u.Age > 120 {
		return fmt.Errorf("age must be between 0 and 150")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
