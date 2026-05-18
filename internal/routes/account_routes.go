package routes

import (
	"github.com/alireza-akbarzadeh/luxe/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupAccountRoutes handle all user acounitng info
func SetupAccountRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	// account endpoints for authenticated users
	protected.GET("/account/summary", ctrl.Account.GetAccountSummary)
}
