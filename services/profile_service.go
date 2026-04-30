package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

// GetUserByID retrieves a user by ID, excluding sensitive fields.
func (s *ProfileService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound(messages.ErrUserNotFound)
		}
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}

// UpdateUserProfile updates non‑sensitive user fields.
func (s *ProfileService) UpdateUserProfile(userID uint, firstName, lastName, phone string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.ErrNotFound(messages.ErrUserNotFound)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &user, nil
}
