package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

// SetupProductRoutes registers all product routes (public + admin protected)
func SetupProductRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// 1. Public product endpoints (no auth)
	public.GET(constants.RouteProducts, ctrl.Product.List)
	public.GET(constants.RouteProducts+"/:identifier", ctrl.Product.GetOne)

	// 2. Admin product endpoints (require JWT + admin role)
	//    No "/admin" prefix – we apply role middleware directly to a sub‑group of the protected router.
	adminProducts := protected.Group(constants.RouteProducts)
	adminProducts.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		adminProducts.POST(constants.RouteRoot, ctrl.Product.Create)
		adminProducts.POST(constants.RouteProductBulk, ctrl.Product.BulkCreate)
		adminProducts.PUT("/:id", ctrl.Product.Update)
		adminProducts.DELETE("/:id", ctrl.Product.Delete)
		adminProducts.DELETE(constants.RouteProductBulk, ctrl.Product.BulkDelete)
	}

	// 3. Admin user management (also protected + admin role)
	adminUsers := protected.Group(constants.RouteUsers)
	adminUsers.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		adminUsers.GET(constants.RouteRoot, ctrl.Profile.GetAllUsers)
	}
}
