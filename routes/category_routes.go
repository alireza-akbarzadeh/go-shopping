package routes

package routes

import (
"github.com/alireza-akbarzadeh/shopping-platform/controllers"
"github.com/alireza-akbarzadeh/shopping-platform/middleware"
"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes registers all category routes (public + admin)
func SetupCategoryRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Public category endpoints (no auth)
	public.GET("/categories", ctrl.Category.List)
	public.GET("/categories/:identifier", ctrl.Category.GetOne)

	// Admin category endpoints (require JWT + admin role)
	admin := protected.Group("/categories")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("/", ctrl.Category.Create)
		admin.POST("/bulk", ctrl.Category.BulkCreate)
		admin.PUT("/:id", ctrl.Category.Update)
		admin.DELETE("/:id", ctrl.Category.Delete)
		admin.DELETE("/bulk", ctrl.Category.BulkDelete)
	}
}
