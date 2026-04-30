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
	authService := services.NewAuthService(db, cfg)
	profileService := services.NewProfileService(db, cfg)

	// 6. Initialize controllers
	ctrl := &controllers.Container{
		Health:  controllers.NewHealthController(db),
		Auth:    controllers.NewAuthController(authService),
		Profile: controllers.NewProfileController(profileService),
	}

	// 7. Setup Gin engine and routes
	engine := setupGin()
	router := routes.NewRouter(engine, ctrl, cfg)
	router.Setup()

	// 8. Start server
	bootStrap(engine, cfg)
}
