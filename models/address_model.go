package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID        uint   `gorm:"not null;index" json:"user_id"`
	AddressType   string `gorm:"not null" json:"address_type" validate:"oneof=shipping billing both"`
	IsDefault     bool   `gorm:"default:false" json:"is_default"`
	RecipientName string `gorm:"not null" json:"recipient_name" validate:"required"`
	Phone         string `gorm:"not null" json:"phone" validate:"required,e164"`
	AddressLine1  string `gorm:"not null" json:"address_line1" validate:"required"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `gorm:"not null" json:"city" validate:"required"`
	State         string `json:"state,omitempty"`
	PostalCode    string `gorm:"not null" json:"postal_code" validate:"required"`
	Country       string `gorm:"not null" json:"country" validate:"required"`
	Instructions  string `json:"instructions,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
