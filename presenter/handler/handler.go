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
	// GetReservation(c echo.Context) error
	// // GetAll(c echo.Context) error
	// SaveReservation(c echo.Context) error
	// // Update(c echo.Context) error
	// // Delete(c echo.Context) error
	// CheckAvailability(c echo.Context) error
	// CheckRoomSchedule(c echo.Context) error
}

type SnackHandler interface {
	GetSnacks(c echo.Context) error
}

type RoomtypeHandler interface {
	GetRoomTypes(c echo.Context) error
	SaveRoomType(c echo.Context) error
	GetRoomType(c echo.Context) error
}

type RoomHandler interface {
	GetRooms(c echo.Context) error
	SaveRoom(c echo.Context) error
	GetRoom(c echo.Context) error
	UpdateRoom(c echo.Context) error
}
