package model

type RoomRequest struct {
	ID         *int   `json:"id,omitempty"`
	Name       string `json:"name" validate:"required,min=2" example:"Aster Room"`
	RoomTypeId int    `json:"room_type_id" validate:"required,min=1" example:"1"`
	Price      int    `json:"price" validate:"required,min=1" example:"10000"`
	Capacity   int    `json:"capacity" validate:"required,min=1" example:"20"`
	ImgUrl     string `json:"imgUrl" validate:"omitempty,url" example:"https://example.com/image.jpg"`
}
