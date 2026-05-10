package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

func SetupCouponRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Public endpoints (list coupons? usually only active ones)
	public.GET("/coupons", ctrl.Coupon.List)

	// Authenticated user endpoints (validate, apply)
	user := protected.Group("/coupons")
	{
		user.POST("/validate", ctrl.Coupon.Validate)
		// POST /apply is part of checkout, not separate
	}

	// Admin endpoints
	admin := protected.Group("/coupons")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("/", ctrl.Coupon.Create)
		admin.PUT("/:id", ctrl.Coupon.Update)
		admin.DELETE("/:id", ctrl.Coupon.Delete)
	}
}
