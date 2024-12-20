package handler

import (
	"context"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type reservationHandler struct {
	useCase usecase.ReservationUseCase
}

func NewReservationHandler(useCase usecase.ReservationUseCase) ReservationHandler {
	return &reservationHandler{useCase: useCase}
}

func (h *reservationHandler) GetReservation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid reservation id"))
	}

	reservation, err := h.useCase.GetReservation(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("reservation not found"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get reservation", reservation))
}

func (h *reservationHandler) SaveReservation(c echo.Context) error {

	reservation := &model.ReservationRequest{}
	err := c.Bind(reservation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request data"))
	}
	savedReservation, err := h.useCase.SaveReservation(context.Background(), reservation)
	if err != nil {
		fmt.Println("Error saving reservation:", err)
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}

	responseData := model.ReservationResponse{
		ReservationID: savedReservation.ID,
		UserID:        savedReservation.UserID,
		RoomID:        savedReservation.RoomID,
		StartTime:     savedReservation.StartTime.Format(time.RFC3339),
		EndTime:       savedReservation.EndTime.Format(time.RFC3339),
		BookingDate:   savedReservation.BookingDate,
		RoomPrice:     savedReservation.RoomPrice,
		SnackPrice:    savedReservation.SnackPrice,
		TotalPrice:    savedReservation.TotalPrice,
		Status:        "booked",
	}

	// Return the structured response
	return c.JSON(http.StatusOK, response.SuccessResponse("success save reservation", responseData))
}

func (h *reservationHandler) CheckAvailability(c echo.Context) error {
	roomId, err := strconv.Atoi(c.QueryParam("room_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}

	startTime := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")

	available, err := h.useCase.CheckAvailability(c.Request().Context(), roomId, startTime, endTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}

	if available {
		return c.JSON(http.StatusOK, response.SuccessResponse("available", nil))
	} else {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("unavailable"))
	}
}

func (h *reservationHandler) CheckRoomSchedule(c echo.Context) error {
	roomId, err := strconv.Atoi(c.QueryParam("room_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}

	date := c.QueryParam("date")

	reservations, err := h.useCase.GetReservationsByRoomAndDate(c.Request().Context(), roomId, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get room schedule", reservations))
}
