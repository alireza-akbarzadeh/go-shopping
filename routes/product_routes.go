package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Public routes – order matters: more specific first
	public.GET("/products", ctrl.Product.List)
	public.GET("/products/:id/related", ctrl.Product.GetRelated)
	public.GET("/products/:id", ctrl.Product.GetOne)

	// Admin routes (protected + admin role)
	adminProducts := protected.Group("/products")
	adminProducts.Use(middleware.RequireRole("admin"))
	{
		adminProducts.POST("/", ctrl.Product.Create)
		adminProducts.POST("/bulk", ctrl.Product.BulkCreate)
		adminProducts.PUT("/:id", ctrl.Product.Update)
		adminProducts.DELETE("/:id", ctrl.Product.Delete)
		adminProducts.DELETE("/bulk", ctrl.Product.BulkDelete)
	}

	adminUsers := protected.Group("/users")
	adminUsers.Use(middleware.RequireRole("admin"))
	{
		adminUsers.GET("/", ctrl.User.GetAllUsers)
	}
}
