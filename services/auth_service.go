package services

import (
	"errors"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// RegisterRequest defines input for registration.
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164"`
}

// Register creates a new user and returns a JWT token.
func (s *AuthService) Register(req RegisterRequest) (string, *models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return "", nil, utils.ErrConflict(messages.ErrEmailAlreadyExists)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         "user",
		IsActive:     true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return "", nil, err
	}

	// Generate JWT token (24h expiration)
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, 24*time.Hour)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// LoginRequest defines input for login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(req LoginRequest) (string, *models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, utils.ErrUnauthorized(messages.ErrInvalidCredentials)
		}
		return "", nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, utils.ErrUnauthorized(messages.ErrAccountDeactivated)
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", nil, utils.ErrUnauthorized(messages.ErrInvalidCredentials)
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.db.Model(&user).Update("last_login_at", now).Error; err != nil {
		// Non‑critical, log but continue
		utils.Log.WithError(err).Warn("failed to update last_login_at")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, 24*time.Hour)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
