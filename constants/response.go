// Package constants defines constant messages (success, error, validation) used across the application.
package constants

const (
	// Auth
	MsgRegistrationSuccess = "registration successful"
	MsgLoginSuccess        = "login successful"

	// Generic
	MsgFetchSuccess   = "data retrieved successfully"
	MsgCreateSuccess  = "resource created successfully"
	MsgUpdateSuccess  = "resource updated successfully"
	MsgDeleteSuccess  = "resource deleted successfully"
	MsgRefreshSuccess = "access token refreshed successfully"
)

var ValidationTagMessages = map[string]string{
	"required": "this field is required",
	"email":    "must be a valid email address",
	"min":      "value is too short",
	"max":      "value is too long",
	"e164":     "phone number must be in E.164 format (e.g., +1234567890)",
}
