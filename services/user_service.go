package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	Role      string `form:"role" binding:"omitempty,oneof=user admin moderator"`
}

type UserServiceInterface interface {
	GetUserByID(userID uint) (*models.User, error)
	GetUsers(filter UserFilter) ([]models.User, int64, error)
	UpdateUserProfile(userID uint, req UpdateProfileRequest) (*models.User, error)
	DeleteUser(userID uint) error
}

type UserService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserService(db *gorm.DB, cfg *config.Config) *UserService {
	return &UserService{db: db, cfg: cfg}
}

// GetUserByID retrieves a user by ID, excluding sensitive fields.
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound()
		}
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}

// UserFilter defines filter and pagination parameters for listing users.
type UserFilter struct {
	Limit     int    `form:"limit" binding:"omitempty,max=100"`
	Offset    int    `form:"offset" binding:"omitempty,min=0"`
	IsActive  *bool  `form:"is_active"`
	Email     string `form:"email"`
	Phone     string `form:"phone"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Role      string `form:"role" binding:"omitempty,oneof=user admin moderator"`
}

// GetUsers retrieve all the users
func (s *UserService) GetUsers(filter UserFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Set defaults
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	query := s.db.Model(&models.User{})

	// Apply filters
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Email != "" {
		query = query.Where("LOWER(email) LIKE LOWER(?)", "%"+filter.Email+"%")
	}
	if filter.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.FirstName != "" {
		query = query.Where("LOWER(first_name) LIKE LOWER(?)", "%"+filter.FirstName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("LOWER(last_name) LIKE LOWER(?)", "%"+filter.LastName+"%")
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	// Count total matching records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	// Stable pagination order
	orderBy := "created_at DESC, id DESC"

	// Fetch paginated results with consistent ordering
	err := query.Limit(filter.Limit).
		Offset(filter.Offset).
		Order(orderBy).
		Find(&users).Error
	if err != nil {
		return nil, 0, utils.ErrInternal(err)
	}

	return users, total, nil
}

// UpdateUserProfile updates non‑sensitive user fields.
func (s *UserService) UpdateUserProfile(userID uint, req UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.ErrNotFound()
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Role = req.Role

	if err := s.db.Save(&user).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}

// DeleteUser removing user with given id
func (s *UserService) DeleteUser(userID uint) error {
	result := s.db.Delete(&models.Product{}, userID)
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("product not found")
	}
	return nil
}
