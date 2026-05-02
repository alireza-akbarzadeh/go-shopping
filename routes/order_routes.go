package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes registers order endpoints (require JWT).
func SetupOrderRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	protected.POST("/orders", ctrl.Order.Checkout)
	protected.GET("/orders/my", ctrl.Order.GetUserOrders)
	protected.GET("/orders/:id", ctrl.Order.GetOrder)

	protected.Use(middleware.RequireRole("admin"))
	{
		protected.GET("/orders", ctrl.Order.ListAllOrders)
	}
}
