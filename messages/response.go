package messages

const (
	MsgRegistrationSuccess = "registration successful"
	MsgLoginSuccess        = "login successful"
	MsgRegistrationFailed  = "registration failed"
	MsgLoginFailed         = "login failed"
	ErrorMissingAuthHeader = "authorization header is missing"
	ErrorInvalidAuthFormat = "authorization header must be in the format: Bearer {token}"
	ErrorUnauthorized      = "You don't have a required permission"
	ErrorForbidden         = "forbidden error"
	ErrorInvalidToken      = "invalid token error"
)
