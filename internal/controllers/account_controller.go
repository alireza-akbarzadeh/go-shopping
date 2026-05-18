package controllers

import (
	"github.com/alireza-akbarzadeh/luxe/internal/constants"
	"github.com/alireza-akbarzadeh/luxe/internal/dto"
	"github.com/alireza-akbarzadeh/luxe/internal/middleware"
	"github.com/alireza-akbarzadeh/luxe/internal/models"
	"github.com/alireza-akbarzadeh/luxe/internal/services"
	"github.com/alireza-akbarzadeh/luxe/internal/utils"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	addressService services.AddressServiceInterface
	likeService    services.UsertLikeServiceInterface
	orderService   services.OrderServiceInterface
	userService    services.UserServiceInterface
}

func NewAccountController(
	addressService services.AddressServiceInterface,
	likeService services.UsertLikeServiceInterface,
	orderService services.OrderServiceInterface,
	userService services.UserServiceInterface,
) *AccountController {
	return &AccountController{
		addressService: addressService,
		likeService:    likeService,
		orderService:   orderService,
		userService:    userService,
	}
}

// GetAccountSummary returns combined user dashboard data.
// @Summary      Get user dashboard summary
// @Description  Returns user profile, default addresses, address count, liked products count, and recent orders (max 3)
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success 200 {object} utils.Response{data=dto.DashboardSummaryResponse}
// @Failure      401 {object} utils.Response
// @Router       /account/summary [get]
func (ac *AccountController) GetAccountSummary(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	// 1. Get user profile (from userService – you need to inject it)
	user, err := ac.userService.GetUserByID(userID)
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch user")
		return
	}

	// 2. Default addresses
	shippingAddr, _ := ac.addressService.GetDefaultAddress(userID, "shipping")
	billingAddr, _ := ac.addressService.GetDefaultAddress(userID, "billing")

	// 3. Address count
	addresses, _ := ac.addressService.List(userID)
	addressCount := len(addresses)

	// 4. Liked products count
	productIDs, _ := ac.likeService.GetUserLikedProductIDs(userID)
	likedCount := len(productIDs)

	// 5. Recent orders
	recentOrders, _, _ := ac.orderService.GetUserOrders(userID, dto.OrderListFilters{Limit: 3, Offset: 0})
	orderDTOs := make([]dto.OrderResponse, len(recentOrders))
	for i, o := range recentOrders {
		orderDTOs[i] = dto.OrderResponse{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			Status:      o.Status,
			TotalAmount: o.TotalAmount,
			CreatedAt:   o.CreatedAt,
		}
	}

	resp := dto.DashboardSummaryResponse{
		ID:                     user.ID,
		Email:                  user.Email,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		Phone:                  user.Phone,
		Role:                   user.Role,
		IsActive:               user.IsActive,
		CreatedAt:              user.CreatedAt,
		DefaultShippingAddress: toAddressDTO(shippingAddr),
		DefaultBillingAddress:  toAddressDTO(billingAddr),
		AddressCount:           addressCount,
		LikedProductsCount:     likedCount,
		RecentOrders:           orderDTOs,
	}
	utils.SuccessResponse(c, "dashboard summary retrieved", resp)
}

// Helper function to convert models.Address to dto.DefaultAddressDTO
func toAddressDTO(addr *models.Address) *dto.DefaultAddressDTO {
	if addr == nil {
		return nil
	}
	return &dto.DefaultAddressDTO{
		ID:           addr.ID,
		AddressLine1: addr.AddressLine1,
		AddressLine2: addr.AddressLine2,
		City:         addr.City,
		State:        addr.State,
		PostalCode:   addr.PostalCode,
		Country:      addr.Country,
		Phone:        addr.Phone,
	}
}
