package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/gin-gonic/gin"
)

func SetupAddressRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	address := protected.Group("/addresses")
	{
		address.POST("/", ctrl.Address.Create)
		address.GET("/", ctrl.Address.List)
		address.PUT("/:id", ctrl.Address.Update)
		address.DELETE("/:id", ctrl.Address.Delete)
		address.PATCH("/:id/default", ctrl.Address.SetDefault)
		address.GET("/default", ctrl.Address.GetDefault)

	}
}
