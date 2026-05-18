package routes

import (
	"github.com/alireza-akbarzadeh/luxe/internal/controllers"
	"github.com/alireza-akbarzadeh/luxe/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupShipmentRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// User shipment endpoints (authenticated)
	public.GET("/shipping-providers", ctrl.Shipment.GetShippingProvider)

	protected.GET("/shipments/:id", ctrl.Shipment.GetShipment)
	protected.GET("/shipments", ctrl.Shipment.GetShipmentsByOrder)

	// Admin shipment endpoints (require admin role)
	admin := protected.Group("/shipments")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("/", ctrl.Shipment.CreateShipment)
		admin.PUT("/:id/status", ctrl.Shipment.UpdateShipmentStatus)
	}
}
