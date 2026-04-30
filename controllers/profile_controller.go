package controllers

import (
	"errors"
	"strconv"

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

// GetAllUsers returns a paginated list of users (admin only).
func (pc *ProfileController) GetAllUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := pc.profileService.GetUsers(limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to fetch users")
		return
	}

	safeUsers := make([]gin.H, len(users))
	for i, u := range users {
		safeUsers[i] = gin.H{
			"id":         u.ID,
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"phone":      u.Phone,
			"role":       u.Role,
			"is_active":  u.IsActive,
			"created_at": u.CreatedAt,
		}
	}

	data := gin.H{
		"users":  safeUsers,
		"limit":  limit,
		"offset": offset,
		"count":  len(users),
	}
	utils.SuccessResponse(c, messages.MsgFetchSuccess, data)
}
