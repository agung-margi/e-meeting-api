package entity

import "time"

type Room struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	RoomTypeId int       `json:"room_type_id"`
	RoomType   string    `json:"room_type"`
	Price      int       `json:"price"`
	Capacity   int       `json:"capacity"`
	ImgUrl     string    `json:"imgUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
