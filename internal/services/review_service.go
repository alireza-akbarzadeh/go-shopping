package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/internal/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/models"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
	"gorm.io/gorm"
)

type ReviewServiceInterface interface {
	Create(userID uint, req dto.CreateReviewRequest) (*models.Review, error)
	Update(userID, reviewID uint, req dto.UpdateReviewRequest) (*models.Review, error)
	Delete(userID, reviewID uint) error
	GetProductReviews(productID uint, limit, offset int) ([]models.Review, int64, error)
	GetUserReviewForProduct(userID, productID uint) (*models.Review, error)
}

type reviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) ReviewServiceInterface {
	return &reviewService{db: db}
}

// Create review and update product rating & count.
func (s *reviewService) Create(userID uint, req dto.CreateReviewRequest) (*models.Review, error) {
	var existing models.Review
	err := s.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existing).Error
	if err == nil {
		return nil, utils.ErrBadRequest("you have already reviewed this product")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrInternal(err)
	}

	// Create review
	review := &models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Title:     req.Title,
	}
	if err := s.db.Create(review).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}

	// Update product rating and reviews count
	s.updateProductStats(req.ProductID)

	return review, nil
}

// Update review and recalc product stats.
func (s *reviewService) Update(userID, reviewID uint, req dto.UpdateReviewRequest) (*models.Review, error) {
	var review models.Review
	err := s.db.Where("id = ? AND user_id = ?", reviewID, userID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("review not found")
		}
		return nil, utils.ErrInternal(err)
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Comment != nil {
		review.Comment = *req.Comment
	}
	if req.Title != nil {
		review.Title = *req.Title
	}

	if err := s.db.Save(&review).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}

	s.updateProductStats(review.ProductID)
	return &review, nil
}

// Delete review and recalc product stats.
func (s *reviewService) Delete(userID, reviewID uint) error {
	var review models.Review
	err := s.db.Where("id = ? AND user_id = ?", reviewID, userID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("review not found")
		}
		return utils.ErrInternal(err)
	}
	if err := s.db.Delete(&review).Error; err != nil {
		return utils.ErrInternal(err)
	}
	s.updateProductStats(review.ProductID)
	return nil
}

// Helper: recalc average rating and count for a product.
func (s *reviewService) updateProductStats(productID uint) {
	var result struct {
		AvgRating float32
		Count     int
	}
	s.db.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Where("product_id = ?", productID).
		Scan(&result)

	s.db.Model(&models.Product{}).Where("id = ?", productID).
		Updates(map[string]interface{}{
			"rating":        result.AvgRating,
			"reviews_count": result.Count,
		})
}

// GetProductReviews returns paginated reviews for a product.
func (s *reviewService) GetProductReviews(productID uint, limit, offset int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{}).Where("product_id = ?", productID)
	query.Count(&total)

	// Preload the User relation – this populates the `User` field
	err := query.Preload("User").Limit(limit).Offset(offset).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

// GetUserReviewForProduct returns the authenticated user's review for a product (if any).
func (s *reviewService) GetUserReviewForProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := s.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // no review, not an error
		}
		return nil, utils.ErrInternal(err)
	}
	return &review, nil
}
