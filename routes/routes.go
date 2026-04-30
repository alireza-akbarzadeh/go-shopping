package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	controllers *controllers.Container
}

func NewRouter(engine *gin.Engine, ctrl *controllers.Container) *Router {
	return &Router{
		engine:      engine,
		controllers: ctrl,
	}
}

func (r *Router) Setup() {
	// Global middleware
	r.engine.Use(middleware.CORS())

	v1 := r.engine.Group("/api/v1")
	{
		// Public routes
		v1.GET("/health", r.controllers.Health.Check)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.controllers.Auth.Register)
			auth.POST("/login", r.controllers.Auth.Login)
		}

		// Protected routes (will add auth middleware later)
		// ...
	}
}

// RegisterMiddleware attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddleware() {
	r.engine.Use(middleware.CORS())
}
