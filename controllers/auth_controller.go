package controllers

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService services.AuthServiceInterface
	validate    *validator.Validate
}

func NewAuthController(authService services.AuthServiceInterface) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles user registration.
// @Summary      Register a new user
// @Description  Create a new account and returns a pair of JWT tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body services.RegisterRequest true "Registration data"
// @Success      201 {object} utils.Response{data=object{access_token=string,refresh_token=string,user=object{id=uint,email=string,first_name=string,last_name=string,role=string}}}
// @Failure      400 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	accessToken, refreshToken, user, err := ctrl.authService.Register(req)
	if err != nil {
		utils.HandleAppError(c, err, constants.MsgRegistrationFailed)
		return
	}

	data := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	}
	utils.CreatedResponse(c, constants.MsgRegistrationSuccess, data)
}

// Login handles user authentication.
// @Summary      Login user
// @Description  Authenticate and return access & refresh tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body services.LoginRequest true "Login credentials"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	accessToken, refreshToken, user, err := ctrl.authService.Login(req)
	if err != nil {
		utils.HandleAppError(c, err, constants.MsgLoginFailed)
		return
	}

	data := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	}
	utils.SuccessResponse(c, constants.MsgLoginSuccess, data)
}

// RefreshRequest represents the body for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh generates a new access token using a valid refresh token.
// @Summary      Refresh access token
// @Description  Obtain a new access token using a valid refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body RefreshRequest true "Refresh token"
// @Success      200 {object} utils.Response{data=object{access_token=string,refresh_token=string}}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/refresh [post]
func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req RefreshRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	newAccessToken, newRefreshToken, err := ctrl.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		utils.HandleAppError(c, err, "failed to refresh tokens")
		return
	}

	utils.SuccessResponse(c, constants.MsgRefreshSuccess, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout revokes the user's refresh token(s).
// @Summary      Logout user
// @Description  Invalidate the refresh token (optional specific token)
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object false "Optional {refresh_token}"
// @Success      200 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	var req services.LogoutRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	if err := ctrl.authService.Logout(userID, req); err != nil {
		utils.InternalServerErrorResponse(c, err, "logout failed")
		return
	}

	utils.SuccessResponse(c, constants.MsgLogoutSuccess, nil)
}

// ChangePassword handles password update for authenticated user.
// @Summary      Change password
// @Description  Change current user's password.
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.ChangePasswordRequest true "Password change request"
// @Success      200 {object} utils.Response{message=string}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /profile/change-password [post]
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
	}
	var req dto.ChangePasswordRequest
	if !utils.BindAndValidate(c, req, ctrl.validate) {
		return
	}
	err := ctrl.authService.ChangePassword(userID, req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to change password")
		return
	}

	utils.SuccessResponse(c, "password changed successfully", nil)
}

func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	_ = ctrl.authService.ForgotPassword(req.Email)
	utils.SuccessResponse(c, "if the email exists, you will receive a reset link", nil)
}

// ResetPassword handles password reset using token.
// @Summary      Reset password
// @Description  Resets password using a valid token received by email.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success      200 {object} utils.Response{message=string}
// @Failure      400 {object} utils.Response
// @Router       /auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	err := ctrl.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		utils.HandleAppError(c, err, "reset password failed")
		return
	}

	utils.SuccessResponse(c, "password reset successfully", nil)
}

// VerifyEmail verifies user's email address with a token.
// @Summary      Verify email
// @Description  Verifies email address using a token sent via email.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        token query string true "Verification token"
// @Success      200 {object} utils.Response{message=string}
// @Failure      400 {object} utils.Response
// @Router       /auth/verify-email [get]
func (ctrl *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.ErrBadRequest("token is required")
		return
	}

	err := ctrl.authService.VerifyEmail(token)
	if err != nil {
		utils.HandleAppError(c, err, "email verification failed")
		return
	}

	utils.SuccessResponse(c, "email verified successfully", nil)
}

// SendVerificationEmail sends a verification email to the authenticated user.
// @Summary      Send verification email
// @Description  Sends an email verification link to the authenticated user's email.
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{message=string}
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /auth/send-verification [post]
func (ctrl *AuthController) SendVerificationEmail(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "user not authenticated")
		return
	}

	err := ctrl.authService.SendVerificationEmail(userID)
	if err != nil {
		utils.HandleAppError(c, err, "failed to send verification email")
		return
	}
	utils.SuccessResponse(c, "verification email sent", nil)
}
