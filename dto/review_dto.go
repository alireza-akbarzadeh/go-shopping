package dto

import "time"

type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment,omitempty"`
	Title     string `json:"title,omitempty"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Title   *string `json:"title,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ProductID  uint      `json:"product_id"`
	UserID     uint      `json:"user_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment,omitempty"`
	IsVerified bool      `json:"is_verified"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
}
