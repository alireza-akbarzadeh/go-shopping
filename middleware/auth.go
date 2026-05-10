package middleware

import (
	"strings"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			// optional: log "token from header"
			return parts[1]
		}
	}
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		// optional: log "token from cookie"
		return cookie
	}
	return ""
}

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
	id, ok := val.(uint)
	return id, ok
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
