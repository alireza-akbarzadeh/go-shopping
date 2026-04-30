package utils

import (
	"github.com/alireza-akbarzadeh/shopping-platform/messages"
)

// AppError represents a custom application error with HTTP status and a message.
type AppError struct {
	Code    int
	Message string
	Err     error // underlying error (optional)
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

// ----- Helper constructors (with optional custom message) -----

// ErrBadRequest returns a 400 Bad Request error.
// If customMsg is provided, it overrides the default message.
func ErrBadRequest(customMsg ...string) *AppError {
	msg := messages.ErrBadRequest
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(400, msg, nil)
}

// ErrNotFound returns a 404 Not Found error.
func ErrNotFound(customMsg ...string) *AppError {
	msg := messages.ErrNotFound
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(404, msg, nil)
}

// ErrUnauthorized returns a 401 Unauthorized error.
func ErrUnauthorized(customMsg ...string) *AppError {
	msg := messages.ErrUnauthorized
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(401, msg, nil)
}

// ErrForbidden returns a 403 Forbidden error.
func ErrForbidden(customMsg ...string) *AppError {
	msg := messages.ErrForbidden
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(403, msg, nil)
}

// ErrConflict returns a 409 Conflict error.
func ErrConflict(customMsg ...string) *AppError {
	msg := messages.ErrConflict
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(409, msg, nil)
}

// ErrValidationFailed returns a 400 Bad Request with validation message.
func ErrValidationFailed(customMsg ...string) *AppError {
	msg := messages.ErrValidationFailed
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(400, msg, nil)
}

// ErrInternal returns a 500 Internal Server Error, always logs the underlying error.
func ErrInternal(err error, customMsg ...string) *AppError {
	msg := messages.ErrInternalServer
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(500, msg, err)
}

// ErrTooManyRequests returns a 429 Too Many Requests error.
func ErrTooManyRequests(customMsg ...string) *AppError {
	msg := messages.ErrTooManyRequests
	if len(customMsg) > 0 && customMsg[0] != "" {
		msg = customMsg[0]
	}
	return NewAppError(429, msg, nil)
}
