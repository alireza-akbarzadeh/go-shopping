package middleware

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token and stores user info in context.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			utils.UnauthorizedResponse(c, constants.ErrorMissingAuthHeader)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, constants.ErrorInvalidToken)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetUserID Optional: helper to get current user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := val.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

// GetUserRole returns the role from context
func GetUserRole(c *gin.Context) (string, bool) {
	val, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := val.(string)
	return role, ok
}

func GuestAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)

		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
