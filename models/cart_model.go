package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Status    string    `gorm:"not null;default:'active'" json:"status"` // active, abandoned, converted
	ExpiresAt time.Time `json:"expires_at"`

	// Associations
	User  User       `gorm:"foreignKey:UserID" json:"-"`
	Items []CartItem `json:"items,omitempty"`
}

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CartID    uint    `gorm:"not null;index" json:"cart_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity" validate:"required,gt=0"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	Color string `gorm:"column:color;size:50;default:''" json:"color"`
	Size  string `gorm:"column:size;size:50;default:''" json:"size"`
	// Associations
	Cart    Cart    `gorm:"foreignKey:CartID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
