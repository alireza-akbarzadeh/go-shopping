package utils

import "github.com/alireza-akbarzadeh/shopping-platform/messages"

// AppError represents a custom application error with status code.
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ErrBadRequest Common app error constructors
func ErrBadRequest(message string) *AppError {
	return NewAppError(400, message, nil)
}

func ErrNotfound() *AppError {
	return NewAppError(404, messages.ErrNotFound, nil)
}

func ErrUnauthorizedApp(message string) *AppError {
	return NewAppError(401, message, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(409, message, nil)
}

func ErrInternal(err error) *AppError {
	return NewAppError(500, messages.ErrInternalServer, err)
}
