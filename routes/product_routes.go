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

	protected.POST("/products/:id/like", ctrl.UserLike.ToggleLike)
	protected.GET("/products/:id/liked", ctrl.UserLike.IsLikedByUser)
	// Admin routes (protected + admin role)
	productActions := protected.Group("/products")
	productActions.Use(middleware.RequireRole("admin"))
	{
		productActions.POST("/", ctrl.Product.Create)
		productActions.POST("/bulk", ctrl.Product.BulkCreate)
		productActions.PUT("/:id", ctrl.Product.Update)
		productActions.DELETE("/:id", ctrl.Product.Delete)
		productActions.DELETE("/bulk", ctrl.Product.BulkDelete)

	}

	adminUsers := protected.Group("/users")
	adminUsers.Use(middleware.RequireRole("admin"))
	{
		adminUsers.GET("/", ctrl.User.GetAllUsers)
	}
}
