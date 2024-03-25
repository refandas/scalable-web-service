package core

import (
	"gorm.io/gorm"
	"time"
)

type SocialMedia struct {
	gorm.Model
	Id             int64      `gorm:"primaryKey"`
	Name           string     `gorm:"not null"`
	SocialMediaUrl string     `gorm:"not null"`
	UserId         int64      `gorm:"not null"`
	CreatedAt      *time.Time `gorm:"autoCreateTime"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime"`
	User           User       `gorm:"foreignKey:UserId"`
}

type SocialMediaRequest struct {
	Name           string `json:"name" validate:"required"`
	SocialMediaUrl string `json:"social_media_url" validate:"required,url"`
	UserId         int64  `validate:"omitempty,gte=1"`
}

type SocialMediaResponse struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	SocialMediaUrl string `json:"social_media_url"`
	UserId         int64  `json:"user_id"`
}

type SocialMediaResponseWithUser struct {
	Id             int64                  `json:"id"`
	Name           string                 `json:"name"`
	SocialMediaUrl string                 `json:"social_media_url"`
	UserId         int64                  `json:"user_id"`
	User           UserResponseSimplified `json:"user"`
}
