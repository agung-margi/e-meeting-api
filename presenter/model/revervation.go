package model

type ReservationRequest struct {
	UserID      int    `json:"userId"`
	RoomID      int    `json:"roomId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	BookingDate string `json:"bookingDate"`
	RoomPrice   int    `json:"roomPrice"`
	SnackPrice  int    `json:"snackPrice"`
	TotalPrice  int    `json:"totalPrice"`
	Status      string `json:"status" default:"booked"`
}

type ReservationResponse struct {
	ReservationID int    `json:"reservationId"`
	UserID        int    `json:"userId"`
	RoomID        int    `json:"roomId"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	BookingDate   string `json:"bookingDate"`
	RoomPrice     int    `json:"roomPrice"`
	SnackPrice    int    `json:"snackPrice"`
	TotalPrice    int    `json:"totalPrice"`
	Status        string `json:"status"`
}

// Name         string `json:"name"`
// Phone        string `json:"phone"`
// Company      string `json:"company"`
// SnackID      int    `json:"snackId"`
// Participants int    `json:"participants"`
