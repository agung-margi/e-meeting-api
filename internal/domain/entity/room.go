package entity

import "time"

type Room struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	RoomType  int       `json:"room_type_id" db:"room_type_id"`
	Price     int       `json:"price" db:"price"`
	Capacity  int       `json:"capacity" db:"capacity"`
	ImgUrl    string    `json:"imgUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoomWithType struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	RoomType  string    `json:"room_type"`
	Price     int       `json:"price"`
	Capacity  int       `json:"capacity"`
	ImgUrl    string    `json:"imgUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
