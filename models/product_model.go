package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a sellable item.
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name              string   `gorm:"not null" json:"name" validate:"required,min=3,max=255"`
	Slug              string   `gorm:"uniqueIndex;not null" json:"slug" validate:"required,slug"`
	Description       string   `json:"description,omitempty"`
	Price             float64  `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	CompareAtPrice    *float64 `gorm:"type:decimal(10,2)" json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	Cost              *float64 `gorm:"type:decimal(10,2)" json:"cost,omitempty" validate:"omitempty,gte=0"`
	SKU               string   `gorm:"uniqueIndex;not null" json:"sku" validate:"required,min=3,max=50"`
	Barcode           string   `json:"barcode,omitempty"`
	Stock             int      `gorm:"not null;default:0" json:"stock" validate:"gte=0"`
	LowStockThreshold int      `gorm:"default:5" json:"low_stock_threshold,omitempty"`
	Weight            *float64 `gorm:"type:decimal(8,2)" json:"weight,omitempty" validate:"omitempty,gte=0"`
	IsDigital         bool     `gorm:"not null;default:false" json:"is_digital"`
	CategoryID        *uint    `json:"category_id,omitempty"`
	Images            []string `gorm:"type:text[]" json:"images,omitempty"`
	Status            string   `gorm:"not null;default:'draft'" json:"status" validate:"oneof=draft active inactive archived"`
	MetaTitle         string   `json:"meta_title,omitempty"`
	MetaDescription   string   `json:"meta_description,omitempty"`
	Rating            float32  `gorm:"type:decimal(3,2);default:0" json:"rating"`
	ReviewsCount      int      `gorm:"default:0" json:"reviews_count"`
	IsNew             bool     `gorm:"default:false" json:"is_new"`
	// Audit
	CreatedByID *uint `json:"created_by,omitempty"`
	UpdatedByID *uint `json:"updated_by,omitempty"`

	// Associations
	Category  *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedBy *User     `gorm:"foreignKey:CreatedByID" json:"-"`
	UpdatedBy *User     `gorm:"foreignKey:UpdatedByID" json:"-"`
}
