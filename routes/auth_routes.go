package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers all authentication routes (public + protected)
func SetupAuthRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// 1. Public auth endpoints (no auth required)
	authGroup := public.Group(constants.RouteAuth)
	{
		authGroup.POST(constants.RouteAuthRegister, ctrl.Auth.Register)
		authGroup.POST(constants.RouteAuthLogin, ctrl.Auth.Login)
	}
	authGroup.POST(constants.RouteAuth+constants.RouteAuthRefresh, ctrl.Auth.Refresh)
	authGroup.POST("/forgot-password", ctrl.Auth.ForgotPassword)
	authGroup.POST("/reset-password", ctrl.Auth.ResetPassword)

	authGroup.GET("/verify-email", ctrl.Auth.VerifyEmail)
	// 2. Protected auth endpoints (require a valid JWT token)
	protected.POST(constants.RouteAuth+"/send-verification", ctrl.Auth.SendVerificationEmail)
	protected.POST(constants.RouteAuth+constants.RouteAuthLogout, ctrl.Auth.Logout)
	protected.POST(constants.RouteAuth+"/change-password", ctrl.Auth.ChangePassword)
}
