package user

import (
	"hashtechy/src/encryption"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
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
		return errors.New(errors.ErrValidation, "name must be between 2 and 50 characters", nil)
	}

	if u.Age < 0 || u.Age > 120 {
		return errors.New(errors.ErrValidation, "age must be between 0 and 150", nil)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New(errors.ErrValidation, "invalid email format", nil)
	}

	logger.Debug("User validation successful for email: %s", u.Email)
	return nil
}

// EncryptEmail encrypts the user's email using base64 encoding
func (u *User) EncryptEmail() (string, error) {
	encryptedEmail, err := encryption.Encrypt(u.Email)
	if err != nil {
		logger.Error("Failed to encrypt email for user: %v", err)
		return "", errors.New(errors.ErrEncryption, "failed to encrypt email", err)
	}
	logger.Debug("Successfully encrypted email for user")
	return encryptedEmail, nil
}

func (u *User) DecryptEmail() error {
	decryptedEmail, err := encryption.Decrypt(u.Email)
	if err != nil {
		logger.Error("Failed to decrypt email for user: %v", err)
		return errors.New(errors.ErrEncryption, "failed to decrypt email", err)
	}
	u.Email = decryptedEmail
	logger.Debug("Successfully decrypted email for user")
	return nil
}
