package entity

import "time"

type Reservation struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	RoomID      int       `json:"roomId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	BookingDate time.Time `json:"bookingDate"`
	RoomPrice   int       `json:"roomPrice"`
	SnackPrice  int       `json:"snackPrice"`
	TotalPrice  int       `json:"totalPrice"`
	Status      string    `json:"status"`
	RoomType    string    `json:"roomType"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ReservationDetails struct {
	ID            int       `json:"id"`
	ReservationID int       `json:"reservationId"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Company       string    `json:"company"`
	SnackID       *int      `json:"snackId"`
	Participants  int       `json:"participants"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type RoomSchedule struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"roomId"`
	RoomName  string    `json:"roomName"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
