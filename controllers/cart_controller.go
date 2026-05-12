package controllers

import (
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CartController struct {
	cartService services.CartServiceInterface
	validate    *validator.Validate
}

func NewCartController(cartService services.CartServiceInterface) *CartController {
	return &CartController{
		cartService: cartService,
		validate:    validator.New(),
	}
}

// AddItem adds a product to the cart.
// @Summary      Add item to cart
// @Description  Add a product to the authenticated user's cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body services.AddItemRequest true "Add item"
// @Success      200 {object} dto.AddItemResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Router       /cart/items [post]
func (ctrl *CartController) AddItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	var req services.AddItemRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	item, err := ctrl.cartService.AddItem(userID, req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to add item")
		return
	}
	resp := dto.AddItemResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "item added to cart",
			Code:    http.StatusOK,
		},
		Data: dto.CartItemData{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		},
	}
	c.JSON(http.StatusOK, resp)
}

// GetCart returns the current user's cart.
// @Summary      Get cart
// @Description  Retrieve all items in the authenticated user's cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.GetCartResponse
// @Failure      401 {object} utils.Response
// @Router       /cart [get]
func (ctrl *CartController) GetCart(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	cart, err := ctrl.cartService.GetCart(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to fetch cart")
		return
	}
	items := make([]dto.CartItemDetail, len(cart.Items))
	var total float64
	for i, item := range cart.Items {
		itemTotal := float64(item.Quantity) * item.Price
		total += itemTotal
		items[i] = dto.CartItemDetail{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Total:     itemTotal,
		}
	}
	resp := dto.GetCartResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgFetchSuccess,
			Code:    http.StatusOK,
		},
		Data: dto.GetCartData{
			Cart: dto.CartData{
				ID:    cart.ID,
				Items: items,
			},
			Total: total,
		},
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateItem updates cart item quantity.
// @Summary      Update cart item quantity
// @Description  Change quantity of a specific cart item
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Param        request body services.UpdateCartItemRequest true "Update quantity"
// @Success      200 {object} dto.EmptyResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Router       /cart/items/{id} [put]
func (ctrl *CartController) UpdateItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid item id")
		return
	}
	var req services.UpdateCartItemRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	if err := ctrl.cartService.UpdateItemQuantity(userID, uint(itemID), req); err != nil {
		utils.HandleAppError(c, err, "failed to update item")
		return
	}

	resp := dto.EmptyResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "cart item updated",
			Code:    http.StatusOK,
		},
	}
	c.JSON(http.StatusOK, resp)
}

// RemoveItem deletes an item from the cart.
// @Summary      Remove cart item
// @Description  Remove a specific item from the cart
// @Tags         Cart
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Success      200 {object} dto.EmptyResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Router       /cart/items/{id} [delete]
func (ctrl *CartController) RemoveItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid item id")
		return
	}

	if err := ctrl.cartService.RemoveItem(userID, uint(itemID)); err != nil {
		utils.HandleAppError(c, err, "failed to remove item")
		return
	}

	resp := dto.EmptyResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "item removed from cart",
			Code:    http.StatusOK,
		},
	}
	c.JSON(http.StatusOK, resp)
}
