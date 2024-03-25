package core

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id              int64  `gorm:"primaryKey"`
	Username        string `gorm:"not null;unique;size:50"`
	Email           string `gorm:"not null;unique;size=150"`
	Password        string `gorm:"not null"`
	Age             uint   `gorm:"not null;check:age >= 8"`
	ProfileImageUrl *string
	CreatedAt       *time.Time `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime"`
}

type UserCreateRequest struct {
	Email           string  `json:"email" validate:"required,email"`
	Username        string  `json:"username" validate:"required,max=50"`
	Age             uint    `json:"age" validate:"required,gte=8"`
	Password        string  `json:"password" validate:"required,min=6"`
	ProfileImageUrl *string `json:"profile_image_url" validate:"omitempty,url"`
}

type UserUpdateRequest struct {
	Email           string  `json:"email" validate:"required,email"`
	Username        string  `json:"username" validate:"required,max=50"`
	Age             uint    `json:"age" validate:"required,gte=8"`
	ProfileImageUrl *string `json:"profile_image_url" validate:"omitempty,url"`
}

type UserResponse struct {
	Id              int64   `json:"id"`
	Email           string  `json:"email"`
	Username        string  `json:"username"`
	Age             uint    `json:"age"`
	ProfileImageUrl *string `json:"profile_image_url"`
}

type UserResponseSimplified struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=150"`
	Password string `json:"password" validate:"required,min=6"`
}
