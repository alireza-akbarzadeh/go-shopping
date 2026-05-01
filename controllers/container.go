// Package controllers
package controllers

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"gorm.io/gorm"
)

type Container struct {
	Health  *HealthController
	Auth    *AuthController
	Profile *ProfileController
	Page    *PageController
	Cart    *CartController
}

// NewContainer initializes all controllers with their dependencies.
func NewContainer(db *gorm.DB, cfg *config.Config, svc *services.Services) *Container {
	return &Container{
		Health:  NewHealthController(db),
		Auth:    NewAuthController(svc.Auth),
		Profile: NewProfileController(svc.Profile),
		Cart:    NewCartController(svc.Cart),
		Page:    NewPageController(),
	}
}
