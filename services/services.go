// Package services contains business logic for app.
package services

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"gorm.io/gorm"
)

type Services struct {
	Auth    AuthServiceInterface
	Profile ProfileServiceInterface
}

func NewServices(db *gorm.DB, cfg *config.Config) *Services {
	return &Services{
		Auth:    NewAuthServices(db, cfg),
		Profile: NewProfileService(db, cfg),
	}
}
