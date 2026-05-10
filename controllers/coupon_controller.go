package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CouponController struct {
	couponService services.CouponServiceInterface
	validate      *validator.Validate
}

func NewCouponController(couponService services.CouponServiceInterface) *CouponController {
	return &CouponController{
		couponService: couponService,
		validate:      validator.New(),
	}
}

// Create a new coupon (admin only).
// @Summary      Create coupon
// @Description  Creates a new discount coupon. Only accessible by users with the "admin" role.
// @Tags         Admin Coupons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCouponRequest true "Coupon creation data"
// @Success      201 {object} utils.Response{data=models.Coupon}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /coupons [post]
func (cc *CouponController) Create(c *gin.Context) {
	var req dto.CreateCouponRequest

	if !utils.BindAndValidate(c, &req, cc.validate) {
		return
	}
	coupon, err := cc.couponService.Create(req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to create coupon")
		return
	}
	utils.CreatedResponse(c, constants.MsgCreateSuccess, coupon)
}

// Update an existing coupon (admin only).
// @Summary      Update coupon
// @Description  Updates a coupon by ID. Only accessible by users with the "admin" role.
// @Tags         Admin Coupons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                       true  "Coupon ID"
// @Param        request body      dto.UpdateCouponRequest   true  "Coupon update data"
// @Success      200     {object}  utils.Response{data=models.Coupon}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      403     {object}  utils.Response
// @Failure      404     {object}  utils.Response
// @Failure      409     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /coupons/{id} [put]
func (cc *CouponController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid coupon id")
		return
	}
	var req dto.UpdateCouponRequest
	if !utils.BindAndValidate(c, &req, cc.validate) {
		return
	}
	coupon, err := cc.couponService.Update(uint(id), req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to update coupon")
		return
	}
	utils.SuccessResponse(c, constants.MsgUpdateSuccess, coupon)
}

// Delete a coupon (admin only).
// @Summary      Delete coupon
// @Description  Soft-deletes a coupon by ID. Only accessible by users with the "admin" role.
// @Tags         Admin Coupons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Coupon ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /coupons/{id} [delete]
func (cc *CouponController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid coupon id")
		return
	}
	err = cc.couponService.Delete(uint(id))
	if err != nil {
		utils.HandleAppError(c, err, "failed to delete coupon")
		return
	}
	utils.SuccessResponse(c, constants.MsgDeleteSuccess, nil)
}

// Validate checks if a coupon code is applicable to the user's cart.
// @Summary      Validate coupon
// @Description  Validates a coupon code for the authenticated user's order total. Returns discount amount and final total if valid.
// @Tags         Coupons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Coupon validation request" SchemaExample({"code":"SAVE10","order_total":100.0})
// @Success      200 {object} utils.Response{data=object{coupon=models.Coupon,discount_amount=float64,final_total=float64}}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /coupons/validate [post]
func (cc *CouponController) Validate(c *gin.Context) {
	var req struct {
		Code       string  `json:"code" validate:"required"`
		OrderTotal float64 `json:"order_total" validate:"required,gt=0"`
	}
	if !utils.BindAndValidate(c, &req, cc.validate) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	coupon, discount, err := cc.couponService.ValidateCoupon(req.Code, userID, req.OrderTotal)
	if err != nil {
		utils.HandleAppError(c, err, "coupon validation failed")
		return
	}
	utils.SuccessResponse(c, "coupon valid", gin.H{
		"coupon":          coupon,
		"discount_amount": discount,
		"final_total":     req.OrderTotal - discount,
	})
}

// List returns a paginated list of coupons with optional filters (admin only).
// @Summary      List coupons
// @Description  Returns a paginated list of coupons with optional filters (admin only)
// @Tags         Admin Coupons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit         query     int     false  "Items per page"  default(20)  minimum(1)  maximum(100)
// @Param        offset        query     int     false  "Offset (skip number of items)"  default(0)  minimum(0)
// @Param        code          query     string  false  "Filter by coupon code (partial match)"
// @Param        is_active     query     bool    false  "Filter by active status"
// @Param        discount_type query     string  false  "Filter by discount type (percentage/fixed)"
// @Param        start_date    query     string  false  "Filter by start date (ISO 8601)"
// @Param        end_date      query     string  false  "Filter by end date (ISO 8601)"
// @Success      200           {object}  utils.Response{data=object{coupons=[]models.Coupon,total=int,limit=int,offset=int}}
// @Failure      400           {object}  utils.Response
// @Failure      401           {object}  utils.Response
// @Failure      403           {object}  utils.Response
// @Failure      500           {object}  utils.Response
// @Router       /coupons [get]
func (cc *CouponController) List(c *gin.Context) {
	var filters dto.CouponListFilters
	if !utils.BindAndValidateQuery(c, &filters, cc.validate) {
		return
	}

	coupons, total, err := cc.couponService.List(filters)
	if err != nil {
		utils.HandleAppError(c, err, "failed to list coupons")
		return
	}

	data := gin.H{
		"coupons": coupons,
		"total":   total,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, data)
}
