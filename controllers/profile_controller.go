package controllers

import (
	"github.com/alireza-akbarzadeh/shopping-platform/messages"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileController struct {
	db *gorm.DB
}

func NewProfileController(db *gorm.DB) *ProfileController {
	return &ProfileController{db: db}
}

// GetProfile returns the authenticated user's profile.
func (pc *ProfileController) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, messages.ErrorUnauthorized)
		return
	}

	var user models.User
	if err := pc.db.First(&user, userID).Error; err != nil {
		utils.NotFoundResponse(c, messages.ErrorUserNotFound)
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
	utils.SuccessResponse(c, messages.SuccessFetch, data)
}
