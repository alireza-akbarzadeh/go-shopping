package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	ProductID  uint   `gorm:"not null;index" json:"product_id"`
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	Rating     int    `gorm:"not null" json:"rating" validate:"min=1,max=5"`
	Comment    string `json:"comment,omitempty"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`

	Product Product `gorm:"foreignKey:ProductID" json:"-"`
	User    User    `gorm:"foreignKey:UserID" json:"-"`
}
