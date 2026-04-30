package controllers

import (
	"errors"
	"net/http"

	"github.com/alireza-akbarzadeh/shopping-platform/messages"
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
// @Description  Create a new account and return access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body services.RegisterRequest true "Registration data"
// @Success      201 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := ctrl.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, formatValidationErrors(err))
		return
	}

	accessToken, refreshToken, user, err := ctrl.authService.Register(req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case http.StatusConflict:
				utils.ConflictResponse(c, appErr.Message)
				return
			case http.StatusBadRequest:
				utils.ErrorResponse(c, http.StatusBadRequest, appErr.Message)
				return
			default:
				utils.InternalServerErrorResponse(c, err, messages.MsgRegistrationFailed)
				return
			}
		}
		utils.InternalServerErrorResponse(c, err, messages.MsgRegistrationFailed)
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
	utils.CreatedResponse(c, messages.MsgRegistrationSuccess, data)
}

// Login handles user authentication.
// @Summary      Login user
// @Description  Authenticate and return access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body services.LoginRequest true "Login credentials"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := ctrl.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, formatValidationErrors(err))
		return
	}

	accessToken, refreshToken, user, err := ctrl.authService.Login(req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case http.StatusUnauthorized:
				utils.UnauthorizedResponse(c, appErr.Message)
				return
			default:
				utils.InternalServerErrorResponse(c, err, messages.MsgLoginFailed)
				return
			}
		}
		utils.InternalServerErrorResponse(c, err, messages.MsgLoginFailed)
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
	utils.SuccessResponse(c, messages.MsgLoginSuccess, data)
}

// RefreshRequest represents the body for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh generates a new access token using a valid refresh token.
// @Summary      Refresh access token
// @Description  Obtain a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh token"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/refresh [post]
func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := ctrl.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, formatValidationErrors(err))
		return
	}

	newAccessToken, newRefreshToken, err := ctrl.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) && appErr.Code == http.StatusUnauthorized {
			utils.UnauthorizedResponse(c, appErr.Message)
			return
		}
		utils.InternalServerErrorResponse(c, err, "failed to refresh tokens")
		return
	}

	utils.SuccessResponse(c, messages.MsgRefreshSuccess, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout revokes the user's refresh token(s).
// @Summary      Logout user
// @Description  Invalidate the refresh token (optional specific token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body object false "Optional {refresh_token}"
// @Success      200 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, messages.ErrUnauthorized)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	_ = c.ShouldBindJSON(&req) // optional

	if err := ctrl.authService.Logout(userID, req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, err, "logout failed")
		return
	}

	utils.SuccessResponse(c, messages.MsgLogoutSuccess, nil)
}
