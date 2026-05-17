package routes

import (
	"github.com/alireza-akbarzadeh/luxe/internal/constants"
	"github.com/alireza-akbarzadeh/luxe/internal/controllers"
	"github.com/alireza-akbarzadeh/luxe/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupShipmentRoutes registers shipment endpoints.
func SetupShipmentRoutes(protected *gin.RouterGroup, ctrl *controllers.Container) {
	// User shipment endpoints
	protected.GET("/shipments/:id", ctrl.Shipment.GetShipment)
	protected.GET("/shipments", ctrl.Shipment.GetShipmentsByOrder)

	// Admin shipment endpoints
	admin := protected.Group(constants.RouteRoot)
	admin.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		admin.POST("/admin/shipments", ctrl.Shipment.CreateShipment)
		admin.PUT("/admin/shipments/:id/status", ctrl.Shipment.UpdateShipmentStatus)
	}
}
