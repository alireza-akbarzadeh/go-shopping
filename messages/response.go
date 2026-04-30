package messages

const (
	// Auth
	MsgRegistrationSuccess = "registration successful"
	MsgLoginSuccess        = "login successful"

	// Generic
	MsgFetchSuccess  = "data retrieved successfully"
	MsgCreateSuccess = "resource created successfully"
	MsgUpdateSuccess = "resource updated successfully"
	MsgDeleteSuccess = "resource deleted successfully"
)

// Optional: Helper to map validation tags to user-friendly messages (used in controllers)
var ValidationTagMessages = map[string]string{
	"required": "this field is required",
	"email":    "must be a valid email address",
	"min":      "value is too short",
	"max":      "value is too long",
	"e164":     "phone number must be in E.164 format (e.g., +1234567890)",
}
