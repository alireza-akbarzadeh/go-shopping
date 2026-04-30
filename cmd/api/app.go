package main

import "github.com/gin-gonic/gin"

// setupGin creates the Gin engine with default middleware.
func setupGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	return engine
}
