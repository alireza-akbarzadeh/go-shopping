// Package routes defines HTTP routing, middleware registration, and endpoint groupings.
package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/alireza-akbarzadeh/shopping-platform/docs" // This line is critical
)

type Router struct {
	engine      *gin.Engine
	controllers *controllers.Container
	cfg         *config.Config
}

func NewRouter(engine *gin.Engine, ctrl *controllers.Container, cfg *config.Config) *Router {
	return &Router{
		engine:      engine,
		controllers: ctrl,
		cfg:         cfg,
	}
}

func (r *Router) Setup() {
	r.engine.Use(middleware.CORS())

	v1 := r.engine.Group("/api/v1")
	r.engine.GET("/", r.controllers.Page.LandingPage)
	r.engine.Static("/static", "./views/static")

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		v1.GET("/health", r.controllers.Health.Check)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.controllers.Auth.Register)
			auth.POST("/login", r.controllers.Auth.Login)
		}

		v1.POST("/auth/refresh", r.controllers.Auth.Refresh)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(r.cfg))
		{
			protected.GET("/profile", r.controllers.Profile.GetProfile)
			protected.PUT("/profile", r.controllers.Profile.UpdateProfile)
			protected.POST("/auth/logout", r.controllers.Auth.Logout)

			// Cart routes
			protected.GET("/cart", r.controllers.Cart.GetCart)
			protected.POST("/cart/items", r.controllers.Cart.AddItem)
			protected.PUT("/cart/items/:id", r.controllers.Cart.UpdateItem)
			protected.DELETE("/cart/items/:id", r.controllers.Cart.RemoveItem)

			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", r.controllers.Profile.GetAllUsers)
			}
		}
	}
}

// RegisterMiddleware attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddleware(cfg *config.Config) {
	r.engine.Use(middleware.CORS())
}
