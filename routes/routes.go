package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
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
		// Public routes (no auth)
		v1.GET("/health", r.controllers.Health.Check)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.controllers.Auth.Register)
			auth.POST("/login", r.controllers.Auth.Login)
		}

		// Protected routes (require JWT)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(r.cfg))
		{
			// Example: get current user profile
			protected.GET("/profile", r.controllers.Profile.GetProfile)

			// Role-specific example: admin only
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/dashboard", r.controllers.Admin.Dashboard) // we'll implement later
			}
		}
	}
}

// RegisterMiddleware attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddleware(cfg *config.Config) {
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.AuthMiddleware())
	r.engine.Use(middleware.RequireRole())
}
