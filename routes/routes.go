package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine           *gin.Engine
	healthController *controllers.HealthController
}

func NewRouter(
	engine *gin.Engine,
	healthController *controllers.HealthController,
) *Router {
	return &Router{
		engine:           engine,
		healthController: healthController,
	}
}

func (r *Router) Setup() {
	// Global middleware (CORS already added in main.go, but we keep a separate file)
	// Health endpoints require no authentication
	v1 := r.engine.Group("/api/v1")
	{
		v1.GET("/health", r.healthController.Check)
	}

	// Public routes (no auth)
	// (will be added later)

	// Protected routes (JWT)
	// (will be added later)
}

// RegisterMiddleware attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddleware() {
	r.engine.Use(middleware.CORS())
}
