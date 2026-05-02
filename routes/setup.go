package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *Router) Setup() {
	r.RegisterMiddlewares()

	// Static files and landing page
	r.engine.GET("/", r.controllers.Page.LandingPage)
	r.engine.Static("/static", "./views/static")
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := r.engine.Group("/api/v1")
	{
		// Health check (public)
		v1.GET("/health", r.controllers.Health.Check)

		// Public router group (no auth)
		public := v1.Group("/")

		// Protected router group (JWT required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(r.cfg))

		// Each module file receives both groups and registers its endpoints
		SetupAuthRoutes(public, protected, r.controllers)
		SetupProductRoutes(public, protected, r.controllers)
		SetupCategoryRoutes(public, protected, r.controllers)
		SetupUserRoutes(protected, r.controllers)
		SetupCartRoutes(protected, r.controllers)
		SetupOrderRoutes(protected, r.controllers)

	}
}
