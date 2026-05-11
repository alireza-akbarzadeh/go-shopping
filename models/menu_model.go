package models

import "time"

type MenuGroup struct {
	ID           uint       `gorm:"primaryKey"`
	Name         string     `gorm:"unique;not null"`
	DisplayOrder int        `gorm:"default:0"`
	Items        []MenuItem `gorm:"foreignKey:GroupID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MenuItem struct {
	ID           uint   `gorm:"primaryKey"`
	GroupID      uint   `gorm:"index;not null"`
	ParentID     *uint  `gorm:"index;default:null"`
	Label        string `gorm:"not null"`
	Href         *string
	Icon         string
	Permission   *string
	DisplayOrder int        `gorm:"default:0"`
	Children     []MenuItem `gorm:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
