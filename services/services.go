// Package services defines the core business logic of the shopping platform. It includes service interfaces and their implementations for authentication, user profiles, shopping carts, products, categories, orders, and shipments. The Services struct aggregates all individual services for easy dependency injection into controllers and other components.
package services

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/tasks"
	"gorm.io/gorm"
)

type Services struct {
	Auth     AuthServiceInterface
	User     UserServiceInterface
	Cart     CartServiceInterface
	Product  ProductServiceInterface
	Category CategoryServiceInterface
	Order    OrderServiceInterface
	Shipment ShipmentServiceInterface
	Hub      *Hub
}

func NewServices(db *gorm.DB, cfg *config.Config, workerPool *tasks.WorkerPool, hub *Hub) *Services {
	return &Services{
		Auth:     NewAuthServices(db, cfg),
		User:     NewUserService(db, cfg),
		Cart:     NewCartService(db),
		Product:  NewProductService(db),
		Category: NewCategoryService(db),
		Order:    NewOrderService(db, hub),
		Shipment: NewShipmentService(db, workerPool),
		Hub:      hub,
	}
}
