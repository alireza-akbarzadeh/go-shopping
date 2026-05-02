package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers all authentication routes (public + protected)
func SetupAuthRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// 1. Public auth endpoints (no auth required)
	authGroup := public.Group("/auth")
	{
		authGroup.POST("/register", ctrl.Auth.Register)
		authGroup.POST("/login", ctrl.Auth.Login)
	}
	public.POST("/auth/refresh", ctrl.Auth.Refresh)

	// 2. Protected auth endpoints (require a valid JWT token)
	protected.POST("/auth/logout", ctrl.Auth.Logout)
}
