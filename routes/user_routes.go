package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes registers user profile endpoints (authenticated users)
// and admin user management endpoints (admin only).
func SetupUserRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Profile endpoints for authenticated users
	protected.GET("/profile", ctrl.Profile.GetProfile)
	protected.PUT("/profile", ctrl.Profile.UpdateProfile)

	// Admin user management (no "/admin" prefix – apply role middleware directly)
	adminUsers := protected.Group("/users")
	adminUsers.Use(middleware.RequireRole("admin"))
	{
		adminUsers.GET("/users", ctrl.Profile.GetAllUsers)
	}
}
