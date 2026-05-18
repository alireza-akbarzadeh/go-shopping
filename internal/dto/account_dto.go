package dto

import (
	"time"
)

type DefaultAddressDTO struct {
	ID           uint   `json:"id"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
}

type DashboardSummaryResponse struct {
	ID                     uint               `json:"id"`
	Email                  string             `json:"email"`
	FirstName              string             `json:"first_name"`
	LastName               string             `json:"last_name"`
	Phone                  string             `json:"phone"`
	Role                   string             `json:"role"`
	IsActive               bool               `json:"is_active"`
	CreatedAt              time.Time          `json:"created_at"`
	DefaultShippingAddress *DefaultAddressDTO `json:"default_shipping_address,omitempty"`
	DefaultBillingAddress  *DefaultAddressDTO `json:"default_billing_address,omitempty"`
	AddressCount           int                `json:"address_count"`
	LikedProductsCount     int                `json:"liked_products_count"`
	RecentOrders           []OrderResponse    `json:"recent_orders"`
}
