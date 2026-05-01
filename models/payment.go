package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Payment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	OrderID         uint           `gorm:"not null;index" json:"order_id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        string         `gorm:"not null;default:'USD'" json:"currency"`
	Method          string         `gorm:"not null" json:"method"`
	Status          string         `gorm:"not null;default:'pending'" json:"status"`
	TransactionID   string         `gorm:"uniqueIndex" json:"transaction_id,omitempty"`
	GatewayResponse datatypes.JSON `gorm:"type:jsonb" json:"gateway_response,omitempty"`

	Order Order `gorm:"foreignKey:OrderID" json:"-"`
	User  User  `gorm:"foreignKey:UserID" json:"-"`
}
