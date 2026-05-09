// Package services defines the core business logic of the shopping platform. It includes service interfaces and their implementations for authentication, user profiles, shopping carts, products, categories, orders, shipments, and real-time notifications. The Services struct aggregates all individual services for easy dependency injection into controllers and other components.
package services

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/tasks"
	"github.com/alireza-akbarzadeh/shopping-platform/websocket"
	"gorm.io/gorm"
)

type Services struct {
	DB           *gorm.DB
	Auth         AuthServiceInterface
	User         UserServiceInterface
	Cart         CartServiceInterface
	Product      ProductServiceInterface
	Category     CategoryServiceInterface
	Order        OrderServiceInterface
	Shipment     ShipmentServiceInterface
	Notification NotificationServiceInterface
	WebSocketHub *websocket.Hub
}

func NewServices(db *gorm.DB, cfg *config.Config, workerPool *tasks.WorkerPool) *Services {
	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize notification service
	notificationSvc := NewNotificationService(db, wsHub)

	return &Services{
		DB:           db,
		Auth:         NewAuthServices(db, cfg),
		User:         NewUserService(db, cfg),
		Cart:         NewCartService(db),
		Product:      NewProductService(db),
		Category:     NewCategoryService(db),
		Order:        NewOrderService(db, notificationSvc),
		Shipment:     NewShipmentService(db, workerPool, notificationSvc),
		Notification: notificationSvc,
		WebSocketHub: wsHub,
	}
}
