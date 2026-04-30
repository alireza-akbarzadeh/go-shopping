package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/routes"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		utils.Log.WithError(err).Fatal("failed to connect to database")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		utils.Log.WithError(err).Fatal("failed to get database connection object")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	utils.Log.Info("database connection established")

	// 5. Initialize controllers (pass db to health controller)
	healthController := controllers.NewHealthController(db)

	// 6. Create Gin engine with default middleware
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 7. Setup routes
	router := routes.NewRouter(engine, healthController)
	router.RegisterMiddleware()
	router.Setup()

	// 8. Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 9. Start server
	go func() {
		utils.Log.Infof("starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.WithError(err).Fatal("failed to start server")
		}
	}()

	// 10. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.Log.WithError(err).Fatal("server forced to shutdown")
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.WithError(err).Error("failed to close database connection")
	}

	utils.Log.Info("server exited properly")
}
