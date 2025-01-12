package model

import "time"

type ReservationRequest struct {
	RoomID       int    `json:"roomId" example:"2"`
	BookingDate  string `json:"bookingDate" example:"2025-02-01"`
	StartTime    string `json:"startTime" example:"08:00"`
	EndTime      string `json:"endTime" example:"09:00"`
	Name         string `json:"name" example:"John Doe"`
	Phone        string `json:"phone" example:"1234567890"`
	Company      string `json:"company" example:"ABC Company"`
	SnackID      *int   `json:"snackId" example:"1"`
	Participants int    `json:"participants" example:"10"`
	Notes        string `json:"notes" example:"Additional notes"`
}

type ReservationResponse struct {
	ReservationID int                          `json:"reservationId"`
	UserID        int                          `json:"userId"`
	RoomID        int                          `json:"roomId"`
	StartTime     time.Time                    `json:"startTime"`
	EndTime       time.Time                    `json:"endTime"`
	BookingDate   time.Time                    `json:"bookingDate"`
	RoomPrice     int                          `json:"roomPrice"`
	SnackPrice    int                          `json:"snackPrice"`
	TotalPrice    int                          `json:"totalPrice"`
	Status        string                       `json:"status" default:"booked"`
	Details       []ReservationDetailsResponse `json:"details"`
}

type ReservationDetailsResponse struct {
	ReservationID int                             `json:"reservationId"`
	UserID        int                             `json:"userId"`
	RoomDetails   []RoomDetailsResponse           `json:"roomDetails"`
	BookDetails   []BookingDetailsResponse        `json:"bookDetails"`
	SnacksDetails []SnackDetailsResponse          `json:"consumption"`
	TotalPrice    []ReservationTotalPriceResponse `json:"totalPrice"`
	Notes         string                          `json:"notes"`
	CreatedAt     time.Time                       `json:"createdAt"`
	UpdatedAt     time.Time                       `json:"updatedAt"`
	ExpiredAt     time.Time                       `json:"expiredAt"`
}

type RoomDetailsResponse struct {
	RoomID    int    `json:"roomId"`
	RoomName  string `json:"roomName"`
	Capacity  int    `json:"capacity"`
	RoomPrice int    `json:"roomPrice"`
}

type BookingDetailsResponse struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Company         string `json:"company"`
	DateReservation string `json:"dateReservation"`
	Duration        int    `json:"duration"`
	Participants    int    `json:"participants"`
}

type SnackDetailsResponse struct {
	SnackID  *int   `json:"snackId"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
}

type ReservationTotalPriceResponse struct {
	RoomName        string `json:"roomName"`
	Duration        int    `json:"duration"`
	RoomPrice       int    `json:"roomPrice"`
	TotalRoomPrice  int    `json:"totalRoomPrice"`
	SnackName       string `json:"snackName"`
	Qty             int    `json:"Qty"`
	TotalSnackPrice int    `json:"totalSnackPrice"`
	TotalPrice      int    `json:"totalPrice"`
}
