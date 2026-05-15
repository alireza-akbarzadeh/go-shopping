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

type UserLikeController struct {
	likeService services.UsertLikeServiceInterface
	validate    *validator.Validate
}

func NewUserLikeController(ls services.UsertLikeServiceInterface) *UserLikeController {
	return &UserLikeController{
		likeService: ls,
		validate:    validator.New(),
	}
}

// ToggleLike toggles a like on a product for the authenticated user.
// @Summary      Toggle product like
// @Description  Like or unlike a product. Send `{"like": true}` to like, `{"like": false}` to unlike.
// @Tags         Product Likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true  "Product ID"
// @Param        request body      object  true  "Toggle action" SchemaExample({"like":true})
// @Success      200     {object}  utils.Response{data=object{liked=bool}}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      404     {object}  utils.Response
// @Router       /products/{id}/like [post]
func (ctrl *UserLikeController) ToggleLike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product id")
		return
	}

	var req struct {
		Like bool `json:"like" validate:"required"`
	}
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	var liked bool
	if req.Like {
		err = ctrl.likeService.Like(userID, uint(productID))
		liked = true
	} else {
		err = ctrl.likeService.Unlike(userID, uint(productID))
		liked = false
	}
	if err != nil {
		utils.HandleAppError(c, err, "failed to toggle like")
		return
	}

	utils.SuccessResponse(c, "success", gin.H{"liked": liked})
}

// IsLikedByUser checks if the current user has liked a specific product.
// @Summary      Check if product is liked
// @Description  Returns whether the authenticated user has liked the given product.
// @Tags         Product Likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response{data=object{liked=bool}}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /products/{id}/liked [get]
func (ctrl *UserLikeController) IsLikedByUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product id")
		return
	}

	liked, err := ctrl.likeService.IsLikedByUser(userID, uint(productID))
	if err != nil {
		utils.HandleAppError(c, err, "failed to check like status")
		return
	}

	utils.SuccessResponse(c, "success", gin.H{"liked": liked})
}

// GetUserLikedProductIDs returns all product IDs liked by the current user.
// @Summary      Get user's liked product IDs
// @Description  Returns a list of product IDs that the authenticated user has liked.
// @Tags         Product Likes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response{data=object{product_ids=[]int}}
// @Failure      401  {object}  utils.Response
// @Router       /users/me/liked-products [get]
func (ctrl *UserLikeController) GetUserLikedProductIDs(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	ids, err := ctrl.likeService.GetUserLikedProductIDs(userID)
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch liked products")
		return
	}

	utils.SuccessResponse(c, "success", gin.H{"product_ids": ids})
}
