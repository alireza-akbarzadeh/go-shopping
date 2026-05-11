// Package services defines the core business logic of the shopping platform.
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
	Coupon       CouponServiceInterface
	Address      AddressServiceInterface
	Menu         UserMenuServicesInterface
}

func NewServices(db *gorm.DB, cfg *config.Config, workerPool *tasks.WorkerPool) *Services {
	// 1. WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 2. Notification service (depends on db + hub)
	notificationSvc := NewNotificationService(db, wsHub)

	// 3. Coupon service (only depends on db)
	couponSvc := NewCouponService(db)

	// 4. Order service (depends on db, notificationSvc, couponSvc)
	orderSvc := NewOrderService(db, notificationSvc, couponSvc)

	// 5. Assemble all services
	return &Services{
		DB:           db,
		Auth:         NewAuthServices(db, cfg),
		User:         NewUserService(db, cfg),
		Cart:         NewCartService(db),
		Product:      NewProductService(db),
		Category:     NewCategoryService(db),
		Address:      NewAddressService(db),
		Menu:         NewMenuService(db),
		Order:        orderSvc,
		Shipment:     NewShipmentService(db, workerPool, notificationSvc),
		Coupon:       couponSvc,
		Notification: notificationSvc,
		WebSocketHub: wsHub,
	}
}
