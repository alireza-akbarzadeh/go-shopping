package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/internal/models"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
	"gorm.io/gorm"
)

type UsertLikeServiceInterface interface {
	Like(userID, productID uint) error
	Unlike(userID, productID uint) error
	IsLikedByUser(userID, productID uint) (bool, error)
	GetUserLikedProductIDs(userID uint) ([]uint, error)
}

type productLikeService struct {
	db *gorm.DB
}

func NewUserLikeService(db *gorm.DB) UsertLikeServiceInterface {
	return &productLikeService{db: db}
}

// Like prorduct service responsible for likeing the product
func (s *productLikeService) Like(userID, productID uint) error {
	// Check if product exists
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("product not found")
		}
		return utils.ErrInternal(err)
	}

	like := &models.ProductLike{
		UserID:    userID,
		ProductID: productID,
	}

	err := s.db.Create(like).Error
	if err == nil {
		return nil
	}

	if isDuplicateKeyError(err) {
		return nil
	}

	return utils.ErrInternal(err)
}

// unlike product services responsible for unliking the product
func (s *productLikeService) Unlike(userID, productID uint) error {
	result := s.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.ProductLike{})
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	return nil
}

// IsLikedByUser check if already product liked by user
func (s *productLikeService) IsLikedByUser(userID, productID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.ProductLike{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	if err != nil {
		return false, utils.ErrInternal(err)
	}
	return count > 0, nil
}

// GetUserLikedProductIDs find product liked by users
func (s *productLikeService) GetUserLikedProductIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := s.db.Model(&models.ProductLike{}).Where("user_id = ?", userID).Pluck("product_id", &ids).Error
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return ids, nil
}
