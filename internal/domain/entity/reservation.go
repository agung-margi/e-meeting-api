package entity

import (
	"time"
)

type Reservation struct {
	ID              int    `json:"id"`
	UserID          int    `json:"userId" `
	RoomID          int    `json:"roomId" binding:"required"`
	StartTime       string `json:"startTime" binding:"required"`
	EndTime         string `json:"endTime" binding:"required"`
	BookingDate     string `json:"bookingDate" binding:"required"`
	RoomPrice       int    `json:"roomPrice" `
	SnackPrice      int    `json:"snackPrice" `
	TotalSnackPrice int    `json:"totalSnackPrice"`
	TotalRoomPrice  int    `json:"totalRoomPrice"`
	TotalPrice      int    `json:"totalPrice"`
	Status          string `json:"status"`
	RoomType        string `json:"roomType" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Company         string `json:"company" binding:"required"`
	SnackID         *int   `json:"snackId"`
	Participants    int    `json:"participants" binding:"required"`
	Notes           string `json:"notes" `
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
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
	ID          int       `json:"id"`
	RoomID      int       `json:"roomId"`
	RoomName    string    `json:"roomName"`
	BookingDate time.Time `json:"BookingDate"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Status      string    `json:"status"`
}
