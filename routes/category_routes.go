package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes registers all category routes (public + admin)
func SetupCategoryRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Public category endpoints (no auth)
	public.GET(constants.RouteCategories, ctrl.Category.List)
	public.GET(constants.RouteCategories+"/:identifier", ctrl.Category.GetOne)

	// Admin category endpoints (require JWT + admin role)
	admin := protected.Group(constants.RouteCategories)
	admin.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		admin.POST(constants.RouteRoot, ctrl.Category.Create)
		admin.POST(constants.RouteProductBulk, ctrl.Category.BulkCreate)
		admin.PUT("/:id", ctrl.Category.Update)
		admin.DELETE("/:id", ctrl.Category.Delete)
		admin.DELETE(constants.RouteProductBulk, ctrl.Category.BulkDelete)
	}
}
