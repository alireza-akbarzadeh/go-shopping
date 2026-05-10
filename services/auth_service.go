package services

import (
	"errors"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthServiceInterface interface {
	Register(req dto.RegisterRequest) (accessToken, refreshToken string, user *models.User, err error)
	Login(req dto.LoginRequest) (accessToken, refreshToken string, user *models.User, err error)
	RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(userID uint, req LogoutRequest) error
	VerifyEmail(token string) error
	ChangePassword(userID uint, req dto.ChangePasswordRequest) error
	ResetPassword(token string, newPassword string) error
	ForgotPassword(email string) error
	SendVerificationEmail(userID uint) error
}
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthServices(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register creates a new user and returns token pair.
func (s *AuthService) Register(req dto.RegisterRequest) (string, string, *models.User, error) {
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
		Role:         constants.RoleUser,
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
func (s *AuthService) Login(req dto.LoginRequest) (string, string, *models.User, error) {
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
	accessToken, err = utils.GenerateToken(user.ID, user.Email, user.Role, user.FirstName, user.LastName, user.Phone, s.cfg.JWT.Secret, 15*time.Minute)
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
func (s *AuthService) Logout(userID uint, req LogoutRequest) error {
	if req.RefreshToken != "" {
		var token models.RefreshToken
		if err := s.db.Where("user_id = ? AND revoked = ?", userID, false).First(&token).Error; err == nil {
			if utils.CheckPasswordHash(req.RefreshToken, token.Token) {
				token.Revoked = true
				return s.db.Save(&token).Error
			}
		}
	}
	return s.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

// ChangePassword use can change password with this services
func (s *AuthService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	var user models.User
	if err := s.db.Select("id", "password_hash").First(&user, userID).Error; err != nil {
		return utils.ErrNotFound("user not found")
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return utils.ErrUnauthorized("current password is incorrect")
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return utils.ErrInternal(err)
	}
	user.PasswordHash = hashed
	if err := s.db.Save(&user).Error; err != nil {
		return utils.ErrInternal(err)
	}

	return nil
}

// ForgotPassword generates a reset token and sends email.
func (s *AuthService) ForgotPassword(email string) error {
	var user models.User
	// Find user by email
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Don't reveal if user exists – return nil (silent)
		return nil
	}

	// Delete any previous unused tokens for this user
	s.db.Where("user_id = ? AND used_at IS NULL", user.ID).Delete(&models.PasswordResetToken{})

	// Generate a secure random token
	token, err := utils.GenerateRandomToken()
	if err != nil {
		return utils.ErrInternal(err)
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	if err := s.db.Create(&resetToken).Error; err != nil {
		return utils.ErrInternal(err)
	}

	// Send email asynchronously
	go utils.SendPasswordResetEmail(user.Email, token)
	return nil
}

// SendVerificationEmail creates a token and sends verification link.
func (s *AuthService) SendVerificationEmail(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.ErrNotFound("user not found")
	}
	if user.EmailVerifiedAt != nil {
		return utils.ErrBadRequest("email already verified")
	}

	// Invalidate previous tokens
	s.db.Where("user_id = ? AND used_at IS NULL", userID).Delete(&models.EmailVerificationToken{})

	token, err := utils.GenerateRandomToken()
	if err != nil {
		return utils.ErrInternal(err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	vt := models.EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	if err := s.db.Create(&vt).Error; err != nil {
		return utils.ErrInternal(err)
	}

	go utils.SendVerificationEmail(user.Email, token) // you'll implement this
	return nil
}

// VerifyEmail marks email as verified.
func (s *AuthService) VerifyEmail(token string) error {
	var vt models.EmailVerificationToken
	err := s.db.Where("token = ? AND used_at IS NULL AND expires_at > ?", token, time.Now()).
		First(&vt).Error
	if err != nil {
		return utils.ErrBadRequest("invalid or expired verification token")
	}

	var user models.User
	if err := s.db.First(&user, vt.UserID).Error; err != nil {
		return utils.ErrInternal(err)
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := s.db.Save(&user).Error; err != nil {
		return utils.ErrInternal(err)
	}

	vt.UsedAt = &now
	s.db.Save(&vt)
	return nil
}

// ResetPassword uses a valid reset token to set a new password.
func (s *AuthService) ResetPassword(token string, newPassword string) error {
	var resetToken models.PasswordResetToken
	now := time.Now()

	err := s.db.Where("token = ? AND used_at IS NULL AND expires_at > ?", token, now).
		First(&resetToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrBadRequest("invalid or expired reset token")
		}
		return utils.ErrInternal(err)
	}

	var user models.User
	if err := s.db.First(&user, resetToken.UserID).Error; err != nil {
		return utils.ErrInternal(err)
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.ErrInternal(err)
	}

	user.PasswordHash = hashed
	if err := s.db.Save(&user).Error; err != nil {
		return utils.ErrInternal(err)
	}

	resetToken.UsedAt = &now
	if err := s.db.Save(&resetToken).Error; err != nil {
		utils.Log.WithError(err).Warn("failed to mark reset token as used")
	}

	s.db.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Update("revoked", true)

	return nil
}
