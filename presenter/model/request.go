package model

import "github.com/go-playground/validator"

type UserRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	ImgUrl   string `json:"imgUrl" validate:"omitempty,url"`
}

func (u *UserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
