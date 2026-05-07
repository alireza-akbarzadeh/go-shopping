package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupProductRoutes registers all product routes (public + admin protected)
func SetupProductRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// 1. Public product endpoints (no auth)
	public.GET("/products", ctrl.Product.List)
	public.GET("/products/:identifier", ctrl.Product.GetOne)

	// 2. Admin product endpoints (require JWT + admin role)
	//    No "/admin" prefix – we apply role middleware directly to a sub‑group of the protected router.
	adminProducts := protected.Group("/products")
	adminProducts.Use(middleware.RequireRole("admin"))
	{
		adminProducts.POST("/", ctrl.Product.Create)
		adminProducts.POST("/bulk", ctrl.Product.BulkCreate)
		adminProducts.PUT("/:id", ctrl.Product.Update)
		adminProducts.DELETE("/:id", ctrl.Product.Delete)
		adminProducts.DELETE("/bulk", ctrl.Product.BulkDelete)
	}

	// 3. Admin user management (also protected + admin role)
	adminUsers := protected.Group("/users")
	adminUsers.Use(middleware.RequireRole("admin"))
	{
		adminUsers.GET("/", ctrl.User.GetAllUsers)
	}
}
