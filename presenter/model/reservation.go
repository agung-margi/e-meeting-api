package model

import "time"

type ReservationRequest struct {
	UserID       int    `json:"userId" binding:"required"`
	RoomID       int    `json:"roomId"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Company      string `json:"company"`
	SnackID      *int   `json:"snackId"`
	Participants int    `json:"participants"`
	Notes        string `json:"notes"`
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
	ReservationID int    `json:"reservationId"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Company       string `json:"company"`
	SnackID       int    `json:"snackId"`
	Participants  int    `json:"participants"`
	Notes         string `json:"notes"`
}
