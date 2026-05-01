package controllers

import (
	"errors"
	"net/http"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
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
				utils.InternalServerErrorResponse(c, err, constants.MsgRegistrationFailed)
				return
			}
		}
		utils.InternalServerErrorResponse(c, err, constants.MsgRegistrationFailed)
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
				utils.InternalServerErrorResponse(c, err, constants.MsgLoginFailed)
				return
			}
		}
		utils.InternalServerErrorResponse(c, err, constants.MsgLoginFailed)
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

	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := ctrl.authService.Logout(userID, req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, err, "logout failed")
		return
	}

	utils.SuccessResponse(c, constants.MsgLogoutSuccess, nil)
}
