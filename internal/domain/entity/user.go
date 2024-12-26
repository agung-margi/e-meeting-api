package entity

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" `
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"isAdmin"`
	ImgUrl    string    `json:"imgUrl"`
	IsActive  bool      `json:"isActive"`
	Language  string    `json:"language"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
