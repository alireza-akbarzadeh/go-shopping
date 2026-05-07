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
	public.POST(constants.RouteAuth+constants.RouteAuthRefresh, ctrl.Auth.Refresh)
	public.POST("/forgot-password", ctrl.Auth.ForgotPassword)
	public.POST("/reset-password", ctrl.Auth.ResetPassword)

	public.GET("/verify-email", ctrl.Auth.VerifyEmail)
	public.GET("/send-verify-email", ctrl.Auth.SendVerificationEmail)
	// 2. Protected auth endpoints (require a valid JWT token)
	protected.POST(constants.RouteAuth+constants.RouteAuthLogout, ctrl.Auth.Logout)
	protected.POST("/change-password", ctrl.Auth.ChangePassword)
}
