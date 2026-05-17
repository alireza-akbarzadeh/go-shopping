// Package routes defines HTTP routing, middleware registration, and endpoint groupings.
package routes

import (
	_ "github.com/alireza-akbarzadeh/shopping-platform/docs"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/config"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/controllers"
	"github.com/gin-gonic/gin"
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
