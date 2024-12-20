package entity

import "time"

type Reservation struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	RoomID      int       `json:"roomId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	BookingDate string    `json:"bookingDate"`
	RoomPrice   int       `json:"roomPrice"`
	SnackPrice  int       `json:"snackPrice"`
	TotalPrice  int       `json:"totalPrice"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}
