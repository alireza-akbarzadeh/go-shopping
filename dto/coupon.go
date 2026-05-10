package dto

import "time"

type CreateCouponRequest struct {
	Code               string    `json:"code" validate:"required,min=3,max=50"`
	Description        string    `json:"description,omitempty"`
	DiscountType       string    `json:"discount_type" validate:"required,oneof=percentage fixed"`
	DiscountValue      float64   `json:"discount_value" validate:"required,gt=0"`
	MinimumOrderAmount float64   `json:"minimum_order_amount"`
	MaxDiscountAmount  *float64  `json:"max_discount_amount,omitempty"`
	UsageLimit         int       `json:"usage_limit"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateCouponRequest struct {
	Code               *string    `json:"code,omitempty"`
	Description        *string    `json:"description,omitempty"`
	DiscountType       *string    `json:"discount_type,omitempty"`
	DiscountValue      *float64   `json:"discount_value,omitempty"`
	MinimumOrderAmount *float64   `json:"minimum_order_amount,omitempty"`
	MaxDiscountAmount  *float64   `json:"max_discount_amount,omitempty"`
	UsageLimit         *int       `json:"usage_limit,omitempty"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsActive           *bool      `json:"is_active,omitempty"`
}

type CouponListFilters struct {
	Limit        int        `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset       int        `form:"offset" validate:"omitempty,min=0"`
	Code         string     `form:"code" validate:"omitempty"`
	IsActive     *bool      `form:"is_active"`
	DiscountType string     `form:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	StartDate    *time.Time `form:"start_date"`
	EndDate      *time.Time `form:"end_date"`
}
