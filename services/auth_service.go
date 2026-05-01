package services

import (
	"errors"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Register(req RegisterRequest) (accessToken, refreshToken string, user *models.User, err error)
	Login(req LoginRequest) (accessToken, refreshToken string, user *models.User, err error)
	RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error) // changed name & returns
	Logout(userID uint, refreshToken string) error
}
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthServices(db *gorm.DB, cfg *config.Config) *AuthService {
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

// LoginRequest defines input for login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register creates a new user and returns token pair.
func (s *AuthService) Register(req RegisterRequest) (string, string, *models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return "", "", nil, utils.ErrConflict(constants.ErrEmailAlreadyExists)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", nil, utils.ErrInternal(err)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", "", nil, utils.ErrInternal(err)
	}

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
		return "", "", nil, utils.ErrInternal(err)
	}

	accessToken, refreshToken, err := s.GenerateTokenPair(user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

// Login authenticates a user and returns token pair.
func (s *AuthService) Login(req LoginRequest) (string, string, *models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, utils.ErrUnauthorized(constants.ErrInvalidCredentials)
		}
		return "", "", nil, utils.ErrInternal(err)
	}

	if !user.IsActive {
		return "", "", nil, utils.ErrUnauthorized(constants.ErrAccountDeactivated)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", "", nil, utils.ErrUnauthorized(constants.ErrInvalidCredentials)
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.db.Model(&user).Update("last_login_at", now).Error; err != nil {
		utils.Log.WithError(err).Warn("failed to update last_login_at")
	}

	accessToken, refreshToken, err := s.GenerateTokenPair(&user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, &user, nil
}

// GenerateTokenPair creates both access and refresh tokens.
func (s *AuthService) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, err error) {
	accessToken, err = utils.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, 15*time.Minute)
	if err != nil {
		return "", "", utils.ErrInternal(err)
	}

	rawRefresh, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", utils.ErrInternal(err)
	}

	hashedRefresh, err := utils.HashPassword(rawRefresh)
	if err != nil {
		return "", "", utils.ErrInternal(err)
	}

	refreshTokenObj := &models.RefreshToken{
		Token:     hashedRefresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	if err := s.db.Create(refreshTokenObj).Error; err != nil {
		return "", "", utils.ErrInternal(err)
	}

	return accessToken, rawRefresh, nil
}

// RefreshTokens validates a refresh token, revokes it, and returns a new token pair.
func (s *AuthService) RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// 1. Find the non‑revoked, non‑expired refresh token
	var storedToken models.RefreshToken
	now := time.Now()
	err = s.db.Where("revoked = ? AND expires_at > ?", false, now).
		Joins("JOIN users ON users.id = refresh_tokens.user_id").
		Where("users.is_active = ?", true).
		First(&storedToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", utils.ErrUnauthorized(constants.ErrInvalidToken)
		}
		return "", "", utils.ErrInternal(err)
	}

	// 2. Verify the provided refresh token matches the stored hash
	if !utils.CheckPasswordHash(refreshToken, storedToken.Token) {
		return "", "", utils.ErrUnauthorized(constants.ErrInvalidToken)
	}

	// 3. Revoke the old token (rotation)
	storedToken.Revoked = true
	if err := s.db.Save(&storedToken).Error; err != nil {
		return "", "", utils.ErrInternal(err)
	}

	// 4. Get the user
	var user models.User
	if err := s.db.First(&user, storedToken.UserID).Error; err != nil {
		return "", "", utils.ErrUnauthorized(constants.ErrUserNotFound)
	}

	// 5. Generate a fresh token pair (access + new refresh)
	newAccess, newRefresh, err := s.GenerateTokenPair(&user)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

// Logout revokes all refresh tokens for a user (or a specific one).
func (s *AuthService) Logout(userID uint, refreshToken string) error {
	if refreshToken != "" {
		var token models.RefreshToken
		if err := s.db.Where("user_id = ? AND revoked = ?", userID, false).First(&token).Error; err == nil {
			if utils.CheckPasswordHash(refreshToken, token.Token) {
				token.Revoked = true
				return s.db.Save(&token).Error
			}
		}
	}
	return s.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
