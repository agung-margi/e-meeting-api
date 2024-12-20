package model

import (
	"github.com/go-playground/validator"
)

type UserRequest struct {
	Username string `json:"username" validate:"required,min=2" example:"user1"`
	Email    string `json:"email" validate:"required,email" example:"user1@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password"`
}
type UpdateUserRequest struct {
	Password string `json:"password" validate:"omitempty,min=6" example:"password"`
	ImgUrl   string `json:"imgUrl" validate:"omitempty,url" example:"https://example.com/image.jpg"`
	IsActive bool   `json:"isActive" validate:"omitempty" example:"true"`
	// UpdatedAt time.Time `json:"updatedAt" validate:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=2" example:"user1"`
	Password string `json:"password" validate:"required,min=6" example:"password"`
}

func (u *UserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UpdateUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
