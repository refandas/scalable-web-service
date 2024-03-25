package core

import (
	"gorm.io/gorm"
	"time"
)

type Photo struct {
	gorm.Model
	Id        int64  `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Caption   *string
	PhotoUrl  string     `gorm:"not null"`
	UserId    int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	User      User       `gorm:"foreignKey:UserId"`
}

type PhotoCreateRequest struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption" validate:"required"`
	PhotoUrl string `json:"photo_url" validate:"required,url"`
	UserId   int64  `json:"user_id"`
}

type PhotoUpdateRequest struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption" validate:"required"`
	PhotoUrl string `json:"photo_url" validate:"required,url"`
}

type PhotoResponse struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   int64  `json:"user_id"`
}

type PhotoResponseWithUser struct {
	Id       int64                  `json:"id"`
	Title    string                 `json:"title"`
	Caption  string                 `json:"caption"`
	PhotoUrl string                 `json:"photo_url"`
	UserId   int64                  `json:"user_id"`
	User     UserResponseSimplified `json:"user"`
}
