package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"log"
	"net/http"
	"strconv"

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

// SaveReservation
// @Summary Save reservation
// @Description Menyimpan reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param reservation body model.ReservationRequest true "Reservation Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations [post]
func (h *reservationHandler) SaveReservation(c echo.Context) error {
	var req model.ReservationRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}

	reservation, err := h.useCase.SaveReservation(c.Request().Context(), &req)
	if err != nil {
		log.Printf("Error saving reservation: %v", err)
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to save reservation"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("reservation saved successfully", reservation))
}
