package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/gin-gonic/gin"
)

// SetupCartRoutes registers all cart endpoints (require JWT authentication).
func SetupCartRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	protected.GET("/cart", ctrl.Cart.GetCart)
	protected.POST("/cart/items", ctrl.Cart.AddItem)
	protected.PUT("/cart/items/:id", ctrl.Cart.UpdateItem)
	protected.DELETE("/cart/items/:id", ctrl.Cart.RemoveItem)
}
