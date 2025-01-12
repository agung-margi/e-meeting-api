package entity

type Inquiry struct {
	ID              int    `json:"id"`
	UserID          int    `json:"userId"`
	RoomID          int    `json:"roomId"`
	RoomName        string `json:"roomName"`
	RoomType        string `json:"roomType"`
	RoomCapacity    int    `json:"roomCapacity"`
	RoomPrice       int    `json:"roomPrice"`
	BookingDate     string `json:"bookingDate"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	Duration        int    `json:"duration"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Company         string `json:"company"`
	Participants    int    `json:"participants"`
	SnackID         *int   `json:"snackId"`
	SnackName       string `json:"snackName"`
	SnackCategory   string `json:"snackCategory"`
	SnackPrice      int    `json:"snackPrice"`
	TotalRoomPrice  int    `json:"totalRoomPrice"`
	TotalSnackPrice int    `json:"totalSnackPrice"`
	TotalPrice      int    `json:"totalPrice"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}
