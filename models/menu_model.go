package models

import (
	"time"
)

type MenuGroup struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"unique;not null" json:"name"`
	DisplayOrder int        `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Items        []MenuItem `gorm:"foreignKey:GroupID" json:"items,omitempty"`
}

type MenuItem struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	GroupID      uint       `gorm:"index;not null" json:"group_id"`
	ParentID     *uint      `gorm:"index;default:null" json:"parent_id,omitempty"`
	Label        string     `gorm:"not null" json:"label"`
	Href         *string    `json:"href,omitempty"`
	Icon         string     `gorm:"not null" json:"icon"`
	Permission   *string    `json:"permission,omitempty"`
	DisplayOrder int        `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Children     []MenuItem `gorm:"-" json:"children,omitempty"`
}
