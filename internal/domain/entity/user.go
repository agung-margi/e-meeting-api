package entity

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" `
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"isAdmin"`
	ImgUrl    string `json:"imgUrl"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
