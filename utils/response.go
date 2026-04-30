package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse sends a 200 OK response with data.
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse sends a 201 Created response with data.
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response with the given status code.
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

// ValidationErrorResponse sends a 400 Bad Request with validation details.
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "validation failed",
		Errors:  errors,
	})
}

// UnauthorizedResponse sends a 401 Unauthorized error.
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse sends a 403 Forbidden error.
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// NotFoundResponse sends a 404 Not Found error.
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// ConflictResponse sends a 409 Conflict error.
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message)
}

// InternalServerErrorResponse sends a 500 Internal Server Error.
func InternalServerErrorResponse(c *gin.Context, err error, message string) {
	if err != nil {
		Log.WithError(err).Error(message)
	} else {
		Log.Error(message)
	}
	ErrorResponse(c, http.StatusInternalServerError, "internal server error")
}
