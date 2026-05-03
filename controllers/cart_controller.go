package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
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
// @Param        request body object true "Add item" SchemaExample({"product_id":1,"quantity":2})
// @Success      200 {object} utils.Response{data=object{cart_item=object{id=uint,product_id=uint,quantity=int,price=float64}}}
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
	var req struct {
		ProductID uint `json:"product_id" validate:"required,gt=0"`
		Quantity  int  `json:"quantity" validate:"required,gt=0"`
	}
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	item, err := ctrl.cartService.AddItem(userID, req.ProductID, req.Quantity)
	if err != nil {
		utils.HandleAppError(c, err, "failed to add item")
		return
	}
	utils.SuccessResponse(c, "item added to cart", gin.H{
		"cart_item": gin.H{
			"id":         item.ID,
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
			"price":      item.Price,
		},
	})
}

// GetCart returns the current user's cart.
// @Summary      Get cart
// @Description  Retrieve all items in the authenticated user's cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=object{cart=object{id=uint,items=[]object{id=uint,product_id=uint,name=string,quantity=int,price=float64,total=float64},total=float64}}}
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
	items := make([]gin.H, len(cart.Items))
	var total float64
	for i, item := range cart.Items {
		itemTotal := float64(item.Quantity) * item.Price
		total += itemTotal
		items[i] = gin.H{
			"id":         item.ID,
			"product_id": item.ProductID,
			"name":       item.Product.Name,
			"quantity":   item.Quantity,
			"price":      item.Price,
			"total":      itemTotal,
		}
	}
	data := gin.H{
		"cart": gin.H{
			"id":    cart.ID,
			"items": items,
		},
		"total": total,
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, data)
}

// UpdateItem updates cart item quantity.
// @Summary      Update cart item quantity
// @Description  Change quantity of a specific cart item
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Param        request body object true "Update quantity" SchemaExample({"quantity":3})
// @Success      200 {object} utils.Response
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
	var req struct {
		Quantity int `json:"quantity" validate:"required,gt=0"`
	}
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	if err := ctrl.cartService.UpdateItemQuantity(userID, uint(itemID), req.Quantity); err != nil {
		utils.HandleAppError(c, err, "failed to update item")
		return
	}

	utils.SuccessResponse(c, "cart item updated", nil)
}

// RemoveItem deletes an item from the cart.
// @Summary      Remove cart item
// @Description  Remove a specific item from the cart
// @Tags         Cart
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Success      200 {object} utils.Response
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

	utils.SuccessResponse(c, "item removed from cart", nil)
}
