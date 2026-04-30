package middleware

import (
	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that allows only users with one of the allowed roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			utils.UnauthorizedResponse(c, messages.ErrorUnauthorized)
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c, messages.ErrorForbidden)
		c.Abort()
	}
}
