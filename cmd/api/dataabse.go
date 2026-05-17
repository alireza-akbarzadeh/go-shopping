package main

import (
	"time"

	"github.com/alireza-akbarzadeh/luxe/internal/config"
	"github.com/alireza-akbarzadeh/luxe/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connectDatabase establishes connection to PostgreSQL with connection pooling.
func connectDatabase(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		utils.Log.WithError(err).Fatal("failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		utils.Log.WithError(err).Fatal("failed to get database connection object")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	utils.Log.Info("database connection established")
	return db
}

// closeDatabase gracefully closes the database connection.
func closeDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		utils.Log.WithError(err).Error("failed to get database connection object")
		return
	}
	if err := sqlDB.Close(); err != nil {
		utils.Log.WithError(err).Error("failed to close database connection")
	}
}
