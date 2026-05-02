package services

import (
	"errors"
	"fmt"

	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	Checkout(userID uint) (*models.Order, error)
	GetUserOrders(userID uint, limit, offset int) ([]models.Order, int64, error)
	GetOrderByID(orderID uint, userID uint) (*models.Order, error)
}

type orderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderServiceInterface {
	return &orderService{db: db}
}

// Checkout converts the user's active cart into an order.
func (s *orderService) Checkout(userID uint) (*models.Order, error) {
	// 1. Get user's active cart with items (preload product)
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
		// 2.1 Calculate total amount and validate stock
		var totalAmount float64
		for _, item := range cart.Items {
			if item.Product.Stock < item.Quantity {
				return utils.ErrBadRequest(fmt.Sprintf("insufficient stock for product: %s", item.Product.Name))
			}
			totalAmount += item.Price * float64(item.Quantity)
		}

		// 2.2 Create order
		order = &models.Order{
			UserID:      userID,
			OrderNumber: generateOrderNumber(userID),
			Status:      "pending",
			TotalAmount: totalAmount,
			Currency:    "USD",
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2.3 Create order items and update product stock
		for _, item := range cart.Items {
			orderItem := &models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
			if err := tx.Create(orderItem).Error; err != nil {
				return err
			}
			// Decrement stock
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&cart).Update("status", "converted").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Preload("Items.Product").Preload("User").First(order, order.ID)
	return order, nil
}

// GetUserOrders returns all orders for a user (paginated).
func (s *orderService) GetUserOrders(userID uint, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	if err := query.Limit(limit).Offset(offset).
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	return orders, total, nil
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