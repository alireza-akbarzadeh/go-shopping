package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	Checkout(userID uint, couponCode string) (*models.Order, error)
	GetUserOrders(userID uint, filters dto.OrderListFilters) ([]models.Order, int64, error)
	GetOrderByID(orderID uint, userID uint) (*models.Order, error)
	GetAllOrders(filters AdminOrderFilters, limit, offset int) ([]models.Order, int64, error)
	UpdateOverdueOrders() error
	UpdateOrderStatus(orderID uint, status string) error
}

type orderService struct {
	db                  *gorm.DB
	notificationService NotificationServiceInterface
	couponService       CouponServiceInterface
}

func NewOrderService(db *gorm.DB, notificationSvc NotificationServiceInterface, couponSvc CouponServiceInterface) OrderServiceInterface {
	return &orderService{
		db:                  db,
		notificationService: notificationSvc,
		couponService:       couponSvc,
	}
}

// Checkout converts the user's active cart into an order.
func (s *orderService) Checkout(userID uint, couponCode string) (*models.Order, error) {
	// 1. Get active cart (same as before)
	var cart models.Cart
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").
		Preload("Items.Product").
		First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrBadRequest("cart is empty")
		}
		return nil, utils.ErrInternal(err)
	}
	if len(cart.Items) == 0 {
		return nil, utils.ErrBadRequest("cart is empty")
	}

	// 2. Start transaction
	var order *models.Order
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 2.1 Calculate subtotal and validate stock
		var subtotal float64
		for _, item := range cart.Items {
			if item.Product.Stock < item.Quantity {
				return utils.ErrBadRequest(fmt.Sprintf("insufficient stock for product: %s", item.Product.Name))
			}
			subtotal += item.Price * float64(item.Quantity)
		}

		// 2.2 Apply coupon if provided
		var discount float64
		var couponID *uint
		if couponCode != "" {
			coupon, disc, err := s.couponService.ValidateCoupon(couponCode, userID, subtotal)
			if err != nil {
				return err // returns AppError already
			}
			discount = disc
			couponID = &coupon.ID
		}

		totalAmount := subtotal - discount
		if totalAmount < 0 {
			totalAmount = 0
		}

		// 2.3 Create order
		order = &models.Order{
			UserID:      userID,
			OrderNumber: generateOrderNumber(userID),
			Status:      constants.OrderStatusPending,
			TotalAmount: totalAmount,
			Currency:    "USD",
		}
		if err := tx.Create(order).Error; err != nil {
			return utils.ErrInternal(err)
		}

		// 2.4 Create order items and update stock
		for _, item := range cart.Items {
			orderItem := &models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
			if err := tx.Create(orderItem).Error; err != nil {
				return utils.ErrInternal(err)
			}
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return utils.ErrInternal(err)
			}
		}

		// 2.5 Mark cart as converted
		if err := tx.Model(&cart).Update("status", "converted").Error; err != nil {
			return utils.ErrInternal(err)
		}

		// 2.6 Apply coupon usage (creates coupon_usage entry and updates used_count)
		if couponCode != "" && couponID != nil {
			if err := s.couponService.ApplyCoupon(userID, order.ID, couponCode, subtotal); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Reload order with preloads
	s.db.Preload("Items.Product").Preload("User").First(order, order.ID)

	// Send notification (async)
	go func() {
		_ = s.notificationService.CreateNotification(
			userID,
			"order_created",
			"Order Placed Successfully",
			fmt.Sprintf("Your order #%s has been placed and is being processed.", order.OrderNumber),
			map[string]interface{}{
				"order_id":     order.ID,
				"order_number": order.OrderNumber,
				"status":       order.Status,
				"total_amount": order.TotalAmount,
				"currency":     order.Currency,
			},
		)
	}()

	return order, nil
}

// GetUserOrders returns all orders for a user (paginated).
func (s *orderService) GetUserOrders(userID uint, filters dto.OrderListFilters) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Set defaults
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	query := s.db.Model(&models.Order{}).Where("user_id = ?", userID)

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.FromDate != nil {
		query = query.Where("created_at >= ?", filters.FromDate)
	}
	if filters.ToDate != nil {
		query = query.Where("created_at <= ?", filters.ToDate)
	}
	if filters.MinAmount != nil {
		query = query.Where("total_amount >= ?", *filters.MinAmount)
	}
	if filters.MaxAmount != nil {
		query = query.Where("total_amount <= ?", *filters.MaxAmount)
	}

	// Count total matching records (efficient)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	// Pagination with ordering – uses indexes
	if err := query.Limit(filters.Limit).Offset(filters.Offset).
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates an order's status and sends real-time notification
func (s *orderService) UpdateOrderStatus(orderID uint, status string) error {
	var order models.Order
	if err := s.db.Preload("User").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("order not found")
		}
		return utils.ErrInternal(err)
	}

	oldStatus := order.Status
	order.Status = status

	if err := s.db.Save(&order).Error; err != nil {
		return utils.ErrInternal(err)
	}

	// Send real-time notification for status change
	go func() {
		title, message := s.getOrderStatusNotificationMessage(status, order.OrderNumber)
		_ = s.notificationService.CreateNotification(
			order.UserID,
			"order_status_update",
			title,
			message,
			map[string]interface{}{
				"order_id":     order.ID,
				"order_number": order.OrderNumber,
				"old_status":   oldStatus,
				"new_status":   status,
				"updated_at":   order.UpdatedAt,
			},
		)
	}()

	return nil
}

// getOrderStatusNotificationMessage returns appropriate title and message for order status
func (s *orderService) getOrderStatusNotificationMessage(status, orderNumber string) (string, string) {
	switch status {
	case constants.OrderStatusPaid:
		return "Payment Confirmed", fmt.Sprintf("Payment for order #%s has been confirmed.", orderNumber)
	case constants.OrderStatusShipped:
		return "Order Shipped", fmt.Sprintf("Your order #%s has been shipped and is on its way!", orderNumber)
	case constants.OrderStatusDelivered:
		return "Order Delivered", fmt.Sprintf("Your order #%s has been delivered successfully.", orderNumber)
	case constants.OrderStatusCancelled:
		return "Order Cancelled", fmt.Sprintf("Your order #%s has been cancelled.", orderNumber)
	case constants.OrderStatusRefunded:
		return "Order Refunded", fmt.Sprintf("Your order #%s has been refunded.", orderNumber)
	default:
		return "Order Update", fmt.Sprintf("Your order #%s status has been updated to %s.", orderNumber, status)
	}
}

// GetOrderByID returns a single order by ID, verifying ownership.
func (s *orderService) GetOrderByID(orderID uint, userID uint) (*models.Order, error) {
	var order models.Order
	err := s.db.Where("id = ? AND user_id = ?", orderID, userID).
		Preload("Items.Product").
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("order not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &order, nil
}

// AdminOrderFilters adds user_id filter
type AdminOrderFilters struct {
	dto.OrderFilters
	UserID *uint `json:"user_id,omitempty"`
}

// GetAllOrders returns all orders (admin only) with advanced filters and pagination.
func (s *orderService) GetAllOrders(filters AdminOrderFilters, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{})

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.FromDate != nil {
		query = query.Where("created_at >= ?", filters.FromDate)
	}
	if filters.ToDate != nil {
		query = query.Where("created_at <= ?", filters.ToDate)
	}
	if filters.MinAmount != nil {
		query = query.Where("total_amount >= ?", *filters.MinAmount)
	}
	if filters.MaxAmount != nil {
		query = query.Where("total_amount <= ?", *filters.MaxAmount)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	// Paginated results with preload
	if err := query.Limit(limit).Offset(offset).
		Preload("Items.Product").
		Preload("User").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	return orders, total, nil
}

// UpdateOverdueOrders marks orders as 'delayed' if they have been 'paid' for more than 7 days.
func (s *orderService) UpdateOverdueOrders() error {
	cutoff := time.Now().Add(-7 * 24 * time.Hour)

	// Find paid orders older than cutoff that are not yet delivered or cancelled
	var orders []models.Order
	err := s.db.Where("status = ? AND updated_at < ?", constants.OrderStatusPaid, cutoff).
		Not("status IN (?)", []string{constants.OrderStatusDelivered, constants.OrderStatusCancelled, constants.OrderStatusRefunded}).
		Find(&orders).Error
	if err != nil {
		return utils.ErrInternal(err)
	}

	if len(orders) == 0 {
		utils.Log.Info("No overdue orders found")
		return nil
	}

	// Mark them as 'delayed'
	for _, order := range orders {
		oldStatus := order.Status
		order.Status = "delayed"
		if err := s.db.Save(&order).Error; err != nil {
			utils.Log.WithError(err).Errorf("Failed to update order %d to delayed", order.ID)
		} else {
			utils.Log.Infof("Order %d marked as delayed", order.ID)

			// Send real-time notification for delayed order
			go func(order models.Order) {
				_ = s.notificationService.CreateNotification(
					order.UserID,
					"order_delayed",
					"Order Delayed",
					fmt.Sprintf("Your order #%s is experiencing a delay. We apologize for the inconvenience.", order.OrderNumber),
					map[string]interface{}{
						"order_id":     order.ID,
						"order_number": order.OrderNumber,
						"old_status":   oldStatus,
						"new_status":   "delayed",
						"updated_at":   order.UpdatedAt,
					},
				)
			}(order)
		}
	}
	return nil
}
