package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(public, protected *gin.RouterGroup, ctrl *controllers.Container) {
	// Public – anyone can see reviews
	public.GET("/reviews", ctrl.Review.GetProductReviews)

	// Protected – authenticated users can manage their own reviews
	protected.POST("/reviews", ctrl.Review.Create)
	protected.PUT("/reviews/:id", ctrl.Review.Update)
	protected.DELETE("/reviews/:id", ctrl.Review.Delete)
}
