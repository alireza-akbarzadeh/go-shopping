package middleware

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that allows only users with one of the allowed roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := GetUserRole(c)
		if !ok {
			utils.UnauthorizedResponse(c, constants.ErrorUnauthorized)
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c, constants.ErrorForbidden)
		c.Abort()
	}
}
