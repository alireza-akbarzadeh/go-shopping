package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes registers order endpoints (require JWT).
func SetupOrderRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	protected.POST(constants.RouteOrders, ctrl.Order.Checkout)
	protected.GET(constants.RouteOrders+constants.RouteOrdersMy, ctrl.Order.GetUserOrders)
	protected.GET(constants.RouteOrders+"/:id", ctrl.Order.GetOrder)

	protected.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		protected.GET(constants.RouteOrders, ctrl.Order.ListAllOrders)
		protected.PUT(constants.RouteOrders+"/:id/status", ctrl.Order.UpdateOrderStatus)
	}
}
