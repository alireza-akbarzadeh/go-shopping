package dto

import "time"

type OrderFilters struct {
	Status    string
	FromDate  *time.Time
	ToDate    *time.Time
	MinAmount *float64
	MaxAmount *float64
}

type OrderListFilters struct {
	Limit     int        `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int        `form:"offset" validate:"omitempty,min=0"`
	Status    string     `form:"status" validate:"omitempty,oneof=pending paid shipped delivered cancelled refunded"`
	FromDate  *time.Time `form:"from_date"`
	ToDate    *time.Time `form:"to_date"`
	MinAmount *float64   `form:"min_amount" validate:"omitempty,gt=0"`
	MaxAmount *float64   `form:"max_amount" validate:"omitempty,gt=0"`
}
