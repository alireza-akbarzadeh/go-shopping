package routes

import (
	"os"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *Router) Setup() {
	r.RegisterMiddlewares()

	// Static files and landing page
	r.engine.GET(constants.RouteRoot, r.controllers.Page.LandingPage)
	r.engine.Static(constants.RouteStatic, "./views/static")
	r.engine.GET(constants.RouteSwagger, ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.engine.GET("/openapi.json", func(c *gin.Context) {
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to read spec"})
			return
		}
		c.Data(200, "application/json", data)
	})
	// API v1 group
	v1 := r.engine.Group(constants.APIVersionV1)
	{
		// Health check (public)
		v1.GET(constants.RouteHealth, r.controllers.Health.Check)

		// Public router group (no auth)
		public := v1.Group(constants.RouteRoot)

		// Protected router group (JWT required)
		protected := v1.Group(constants.RouteRoot)
		protected.Use(middleware.AuthMiddleware(r.cfg))

		// Each module file receives both groups and registers its endpoints
		SetupAuthRoutes(public, protected, r.controllers)
		SetupProductRoutes(public, protected, r.controllers)
		SetupCategoryRoutes(public, protected, r.controllers)
		SetupUserRoutes(protected, r.controllers)
		SetupCartRoutes(protected, r.controllers)
		SetupOrderRoutes(protected, r.controllers)
		SetupShipmentRoutes(protected, r.controllers)
	}
}
