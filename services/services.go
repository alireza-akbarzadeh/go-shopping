// Package services contains business logic for app.
package services

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/tasks"
	"gorm.io/gorm"
)

type Services struct {
	Auth     AuthServiceInterface
	Profile  ProfileServiceInterface
	Cart     CartServiceInterface
	Product  ProductServiceInterface
	Category CategoryServiceInterface
	Order    OrderServiceInterface
	Shipment ShipmentServiceInterface
}

func NewServices(db *gorm.DB, cfg *config.Config, workerPool *tasks.WorkerPool) *Services {
	return &Services{
		Auth:     NewAuthServices(db, cfg),
		Profile:  NewProfileService(db, cfg),
		Cart:     NewCartService(db),
		Product:  NewProductService(db),
		Category: NewCategoryService(db),
		Order:    NewOrderService(db),
		Shipment: NewShipmentService(db, workerPool),
	}
}
