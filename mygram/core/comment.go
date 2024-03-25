package core

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	Id        int64 `gorm:"primaryKey"`
	UserId    int64 `gorm:"not null"`
	PhotoId   int64 `gorm:"not null"`
	Message   string
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	User      User       `gorm:"foreignKey:UserId"`
	Photo     Photo      `gorm:"foreignKey:PhotoId"`
}

type CommentCreateRequest struct {
	UserId  int64  `validate:"omitempty,gte=1"`
	PhotoId int64  `json:"photo_id" validate:"required,gte=1"`
	Message string `json:"message" validate:"required"`
}

type CommentUpdateRequest struct {
	PhotoId int64  `validate:"omitempty,gte=1"`
	Message string `json:"message" validate:"required"`
}

type CommentResponse struct {
	Id      int64  `json:"id"`
	UserId  int64  `json:"user_id"`
	PhotoId int64  `json:"photo_id"`
	Message string `json:"message"`
}

type CommentResponseWithUserAndPhoto struct {
	Id      int64                  `json:"id"`
	UserId  int64                  `json:"user_id"`
	PhotoId int64                  `json:"photo_id"`
	Message string                 `json:"message"`
	User    UserResponseSimplified `json:"user"`
	Photo   PhotoResponse          `json:"photo"`
}
