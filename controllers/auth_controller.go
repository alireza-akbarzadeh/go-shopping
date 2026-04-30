package controllers

import (
	"errors"
	"net/http"

	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService *services.AuthService
	validate    *validator.Validate
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles user registration.
// @Summary      Register a new user
// @Description  Create a new account and return a JWT token
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
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	token, user, err := ctrl.authService.Register(req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case http.StatusConflict: // 409
				utils.ConflictResponse(c, appErr.Message)
				return
			case http.StatusBadRequest: // 400
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
		"token": token,
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
// @Description  Authenticate and return a JWT token
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
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	token, user, err := ctrl.authService.Login(req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case http.StatusUnauthorized: // 401
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
		"token": token,
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
