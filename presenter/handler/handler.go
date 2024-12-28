package handler

import (
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	GetUser(c echo.Context) error
	SaveUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type ReservationHandler interface {
	GetReservation(c echo.Context) error
	GetReservations(c echo.Context) error
	SaveReservation(c echo.Context) error
	CancelReservation(c echo.Context) error
	PayReservation(c echo.Context) error
	GetRoomSchedule(c echo.Context) error
}

type SnackHandler interface {
	GetSnacks(c echo.Context) error
}

type RoomtypeHandler interface {
	GetRoomTypes(c echo.Context) error
	GetRoomType(c echo.Context) error
}

type RoomHandler interface {
	GetRooms(c echo.Context) error
	SaveRoom(c echo.Context) error
	GetRoom(c echo.Context) error
	UpdateRoom(c echo.Context) error
	DeleteRoom(c echo.Context) error
}

type InquiryHandler interface {
	SaveInquiry(c echo.Context) error
}
