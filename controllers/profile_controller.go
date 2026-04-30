package controllers

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProfileController struct {
	profileService *services.ProfileService
	validate       *validator.Validate
}

func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
		validate:       validator.New(),
	}
}

// GetProfile returns the authenticated user's profile.
func (pc *ProfileController) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, messages.ErrUnauthorized)
		return
	}

	user, err := pc.profileService.GetUserByID(userID)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) && appErr.Code == 404 {
			utils.NotFoundResponse(c, messages.ErrUserNotFound)
			return
		}
		utils.InternalServerErrorResponse(c, err, messages.ErrInternalServer)
		return
	}

	// Return safe user data (exclude password hash)
	data := gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	}
	utils.SuccessResponse(c, messages.MsgFetchSuccess, data)
}

// UpdateProfile updates the authenticated user's profile.
func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" validate:"required,min=1,max=100"`
		LastName  string `json:"last_name" validate:"required,min=1,max=100"`
		Phone     string `json:"phone" validate:"omitempty,e164"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := pc.validate.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, messages.ErrUnauthorized)
		return
	}

	user, err := pc.profileService.UpdateUserProfile(userID, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) && appErr.Code == 404 {
			utils.NotFoundResponse(c, messages.ErrUserNotFound)
			return
		}
		utils.InternalServerErrorResponse(c, err, messages.ErrInternalServer)
		return
	}

	data := gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
	}
	utils.SuccessResponse(c, messages.MsgUpdateSuccess, data)
}
