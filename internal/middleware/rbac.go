package middleware

import (
	"strings"

	"github.com/alireza-akbarzadeh/luxe/internal/constants"
	"github.com/alireza-akbarzadeh/luxe/internal/utils"
	"github.com/gin-gonic/gin"
)

// FIXME: find out the problem
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip role check for user‑facing endpoints (no admin needed)
		skipPaths := []string{
			"/api/v1/addresses",
		}
		for _, p := range skipPaths {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		// OPTIONS preflight should never require a role
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Original role check
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
