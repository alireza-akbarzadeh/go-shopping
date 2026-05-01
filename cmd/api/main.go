package main

import (
	"fmt"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/routes"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

// @title           Shopping Platform API
// @version         1.0
// @description     Production-grade e-commerce backend
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@shopping-platform.com

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Initialize logger
	if err := utils.InitLogger(cfg.Log.Level); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}

	// 3. Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// 4. Connect to database
	db := connectDatabase(cfg)
	defer closeDatabase(db)

	// 5. Initialize services
	services := services.NewServices(db, cfg)

	// 6. Initialize controllers
	ctrl := controllers.NewContainer(db, cfg, services)
	// 7. Setup Gin engine and routes
	engine := setupGin()
	router := routes.NewRouter(engine, ctrl, cfg)
	router.Setup()

	// 8. Start server
	bootStrap(engine, cfg)
}
