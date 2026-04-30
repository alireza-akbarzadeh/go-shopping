package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type ProfileService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewProfileService(db *gorm.DB, cfg *config.Config) *ProfileService {
	return &ProfileService{db: db, cfg: cfg}
}

// GetUserByID retrieves a user by ID, excluding sensitive fields.
func (s *ProfileService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound()
		}
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}

func (s *ProfileService) GetUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := s.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return users, nil
}

// UpdateUserProfile updates non‑sensitive user fields.
func (s *ProfileService) UpdateUserProfile(userID uint, firstName, lastName, phone string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.ErrNotFound()
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}
