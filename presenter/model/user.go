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
	Username string `json:"username" validate:"omitempty,min=2" example:"user1"`
	Email    string `json:"email" validate:"omitempty,email" example:"user1@example.com"`
	Password string `json:"password" validate:"omitempty,min=6" example:"password"`
	IsAdmin  bool   `json:"isAdmin" validate:"omitempty" example:"true"`
	ImgUrl   string `json:"imgUrl" validate:"omitempty,url" example:"https://example.com/image.jpg"`
	IsActive bool   `json:"isActive" validate:"omitempty" example:"true"`
	Language string `json:"language" validate:"omitempty" example:"english"`
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

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}
