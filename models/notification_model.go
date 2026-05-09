package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification represents real-time notifications sent to users
type Notification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"not null;index" json:"type"` // order_update, payment_success, etc.
	Title     string         `gorm:"not null" json:"title"`
	Message   string         `gorm:"not null" json:"message"`
	Data      string         `gorm:"type:json" json:"data,omitempty"` // JSON data for frontend
	IsRead    bool           `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Message represents chat messages between users and support
type Message struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	SenderID   uint           `gorm:"not null;index" json:"sender_id"`
	ReceiverID uint           `gorm:"index" json:"receiver_id,omitempty"` // null for broadcast messages
	RoomID     string         `gorm:"not null;index" json:"room_id"`      // chat room identifier
	Content    string         `gorm:"not null" json:"content"`
	Type       string         `gorm:"default:'text'" json:"type"` // text, image, file
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ChatRoom represents chat conversation rooms
type ChatRoom struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	RoomID        string         `gorm:"uniqueIndex;not null" json:"room_id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	AdminID       *uint          `gorm:"index" json:"admin_id,omitempty"`
	Title         string         `gorm:"not null" json:"title"`
	Status        string         `gorm:"default:'active'" json:"status"` // active, closed
	LastMessageAt time.Time      `json:"last_message_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
