package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"hashtechy/src/errors"
	"hashtechy/src/logger"
	"io"
	"os"

	"github.com/joho/godotenv"
)

var secretKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Error("Failed to load .env file: %v", err)
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "your-32-byte-secret-key-here!!" // Exactly 32 bytes
		logger.Debug("Using default encryption key")
	}

	if len(key) > 32 {
		key = key[:32]
		logger.Warn("Encryption key truncated to 32 bytes")
	} else if len(key) < 32 {
		padding := make([]byte, 32-len(key))
		key = key + string(padding)
		logger.Warn("Encryption key padded to 32 bytes")
	}

	secretKey = []byte(key)
	logger.Info("Encryption initialized successfully")
}

func Encrypt(text string) (string, error) {
	if text == "" {
		return "", errors.New(errors.ErrValidation, "empty text provided", nil)
	}

	plaintext := []byte(text)
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", errors.New(errors.ErrEncryption, "failed to create cipher", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New(errors.ErrEncryption, "failed to generate IV", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	encoded := base64.URLEncoding.EncodeToString(ciphertext)
	logger.Debug("Successfully encrypted text of length %d", len(text))
	return encoded, nil
}

func Decrypt(cryptoText string) (string, error) {
	if cryptoText == "" {
		return "", errors.New(errors.ErrValidation, "empty crypto text provided", nil)
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", errors.New(errors.ErrEncryption, "failed to decode base64 string", err)
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", errors.New(errors.ErrEncryption, "failed to create cipher", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New(errors.ErrValidation, "ciphertext too short", nil)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	logger.Debug("Successfully decrypted text of length %d", len(ciphertext))
	return string(ciphertext), nil
}
