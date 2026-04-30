package messages

const (
	EMAIL_ALREADY_REGISTERED = "email already registered"
	// Auth errors
	ErrEmailAlreadyExists = "email already registered"
	ErrInvalidCredentials = "invalid email or password"
	ErrAccountDeactivated = "account is deactivated"
	ErrUserNotFound       = "user not found"
	ErrInvalidToken       = "invalid or expired token"

	// Validation errors
	ErrValidationFailed     = "validation failed"
	ErrInvalidRequestFormat = "invalid request format"

	// Generic errors
	ErrInternalServer = "internal server error"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "access denied"
	ErrNotFound       = "resource not found"

	// Success messages
	MsgLogoutSuccess = "logout successful"
)
