package services

import (
	"errors"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type CartServiceInterface interface {
	GetOrCreateCart(userID uint) (*models.Cart, error)
	AddItem(userID uint, productID uint, quantity int) (*models.CartItem, error)
	UpdateItemQuantity(userID uint, cartItemID uint, quantity int) error
	RemoveItem(userID uint, cartItemID uint) error
	GetCart(userID uint) (*models.Cart, error)
	ClearCart(userID uint) error
}

type cartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) CartServiceInterface {
	return &cartService{db: db}
}

// GetOrCreateCart returns existing active cart or creates a new one.
func (s *cartService) GetOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").
		Preload("Items.Product"). // optional: preload for view
		First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrInternal(err)
	}

	// Create new cart
	cart = models.Cart{
		UserID:    userID,
		Status:    "active",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.db.Create(&cart).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &cart, nil
}

// AddItem adds a product to the cart.
func (s *cartService) AddItem(userID uint, productID uint, quantity int) (*models.CartItem, error) {
	if quantity <= 0 {
		return nil, utils.ErrBadRequest("quantity must be positive")
	}

	// Get product and check stock
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product not found")
		}
		return nil, utils.ErrInternal(err)
	}
	if product.Status != "active" {
		return nil, utils.ErrBadRequest("product is not available")
	}
	if product.Stock < quantity {
		return nil, utils.ErrBadRequest("insufficient stock")
	}

	// Get or create cart
	cart, err := s.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists
	var cartItem models.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error
	if err == nil {
		// Update quantity
		newQty := cartItem.Quantity + quantity
		if product.Stock < newQty {
			return nil, utils.ErrBadRequest("insufficient stock for updated quantity")
		}
		cartItem.Quantity = newQty
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, utils.ErrInternal(err)
		}
		return &cartItem, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrInternal(err)
	}

	// Create new cart item with price snapshot
	cartItem = models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
	}
	if err := s.db.Create(&cartItem).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &cartItem, nil
}

// UpdateItemQuantity modifies existing cart item quantity.
func (s *cartService) UpdateItemQuantity(userID uint, cartItemID uint, quantity int) error {
	if quantity <= 0 {
		return utils.ErrBadRequest("quantity must be positive")
	}

	var cartItem models.CartItem
	if err := s.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ? AND carts.status = ?", cartItemID, userID, "active").
		First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("cart item not found")
		}
		return utils.ErrInternal(err)
	}

	// Validate stock
	var product models.Product
	if err := s.db.First(&product, cartItem.ProductID).Error; err != nil {
		return utils.ErrInternal(err)
	}
	if product.Stock < quantity {
		return utils.ErrBadRequest("insufficient stock")
	}

	cartItem.Quantity = quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

// RemoveItem deletes a cart item.
func (s *cartService) RemoveItem(userID uint, cartItemID uint) error {
	result := s.db.Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ? AND status = ?)",
		cartItemID, userID, "active").Delete(&models.CartItem{})
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("cart item not found")
	}
	return nil
}

// GetCart returns full cart with items for the user.
func (s *cartService) GetCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").
		Preload("Items.Product").
		First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty cart (not created yet)
			return &models.Cart{UserID: userID, Items: []models.CartItem{}}, nil
		}
		return nil, utils.ErrInternal(err)
	}
	return &cart, nil
}

// ClearCart removes all items from active cart.
func (s *cartService) ClearCart(userID uint) error {
	// Delete all cart items belonging to user's active cart
	err := s.db.Where("cart_id IN (SELECT id FROM carts WHERE user_id = ? AND status = ?)", userID, "active").
		Delete(&models.CartItem{}).Error
	return err
}
