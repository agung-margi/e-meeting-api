package entity

import "time"

type Room struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	RoomType  int       `json:"room_type_id"`
	Price     int       `json:"price"`
	Capacity  int       `json:"capacity"`
	ImgUrl    string    `json:"imgUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoomWithType struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	RoomType  string    `json:"room_type"` // Room type as a string
	Price     int       `json:"price"`
	Capacity  int       `json:"capacity"`
	ImgUrl    string    `json:"imgUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
