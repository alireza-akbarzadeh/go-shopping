package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/internal/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/services"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReviewController struct {
	reviewService services.ReviewServiceInterface
	validate      *validator.Validate
}

func NewReviewController(svc services.ReviewServiceInterface) *ReviewController {
	return &ReviewController{
		reviewService: svc,
		validate:      validator.New(),
	}
}

// Create a review
// @Summary      Create product review
// @Description  Leave a rating and comment for a product
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateReviewRequest true "Review data"
// @Success      201 {object} utils.Response{data=models.Review}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Router       /reviews [post]
func (rc *ReviewController) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	var req dto.CreateReviewRequest
	if !utils.BindAndValidate(c, &req, rc.validate) {
		return
	}
	review, err := rc.reviewService.Create(userID, req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to create review")
		return
	}
	utils.CreatedResponse(c, "review submitted", review)
}

// Update a review
// @Summary      Update a review
// @Description  Modify rating or comment of an existing review
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                       true  "Review ID"
// @Param        request body      dto.UpdateReviewRequest  true  "Updated review data"
// @Success      200     {object}  utils.Response{data=models.Review}
// @Router       /reviews/{id} [put]
func (rc *ReviewController) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid review id")
		return
	}
	var req dto.UpdateReviewRequest
	if !utils.BindAndValidate(c, &req, rc.validate) {
		return
	}
	review, err := rc.reviewService.Update(userID, uint(id), req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to update review")
		return
	}
	utils.SuccessResponse(c, "review updated", review)
}

// Delete a review
// @Summary      Delete a review
// @Description  Remove a review by ID
// @Tags         Reviews
// @Security     BearerAuth
// @Param        id   path      int  true  "Review ID"
// @Success      200  {object}  utils.Response
// @Router       /reviews/{id} [delete]
func (rc *ReviewController) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid review id")
		return
	}
	err = rc.reviewService.Delete(userID, uint(id))
	if err != nil {
		utils.HandleAppError(c, err, "failed to delete review")
		return
	}
	utils.SuccessResponse(c, "review deleted", nil)
}

// GetProductReviews product reviews (public)
// @Summary      Get product reviews
// @Description  Returns paginated reviews for a specific product
// @Tags         Reviews
// @Param        product_id query int true "Product ID"
// @Param        limit      query int false "Items per page" default(20)
// @Param        offset     query int false "Offset" default(0)
// @Success      200 {object} utils.Response{data=object{reviews=[]dto.ReviewResponse,total=int,limit=int,offset=int}}
// @Router       /reviews [get]
func (rc *ReviewController) GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product_id")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 {
		limit = 2
	}
	if limit > 100 {
		limit = 100
	}
	reviews, total, err := rc.reviewService.GetProductReviews(uint(productID), limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to fetch reviews")
		return
	}

	// Build response slice with author names
	responseReviews := make([]dto.ReviewResponse, len(reviews))
	for i, rev := range reviews {
		author := ""
		// rev.User is now populated because we used Preload("User")
		if rev.User.ID != 0 { // user exists
			author = rev.User.FirstName + " " + rev.User.LastName
		}
		responseReviews[i] = dto.ReviewResponse{
			ID:         rev.ID,
			CreatedAt:  rev.CreatedAt,
			UpdatedAt:  rev.UpdatedAt,
			ProductID:  rev.ProductID,
			UserID:     rev.UserID,
			Rating:     rev.Rating,
			Comment:    rev.Comment,
			IsVerified: rev.IsVerified,
			Title:      rev.Title,
			Author:     author,
		}
	}

	data := gin.H{
		"reviews": responseReviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, data)
}
