package errors

import (
	"fmt"
	"hashtechy/src/logger"
)

type ErrorCode string

const (
	ErrDatabase     ErrorCode = "DATABASE_ERROR"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrEncryption   ErrorCode = "ENCRYPTION_ERROR"
	ErrNetwork      ErrorCode = "NETWORK_ERROR"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code ErrorCode, message string, err error) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
	logger.Error("%v", appErr)
	return appErr
}
