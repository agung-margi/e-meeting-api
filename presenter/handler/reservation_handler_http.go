package handler

import (
	"e-meeting-api/internal/domain/entity"
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

// GetReservation
// @Summary Get reservation by id
// @Description Mendapatkan reservation berdasarkan id
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /reservations/{id} [get]
func (h *reservationHandler) GetReservation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid reservation id"))
	}

	userIDinterface := c.Get("user_id")
	userIDfloat, ok := userIDinterface.(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("invalid token claims"))
	}
	userID := int(userIDfloat)
	isAdmin := c.Get("is_admin").(bool)

	reservation, err := h.useCase.GetReservation(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("reservation not found"))
	}

	if reservation.UserID != userID && !isAdmin {
		return c.JSON(http.StatusForbidden, response.ErrorResponse("forbidden", http.StatusForbidden))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get reservation", reservation))
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

// GetReservations
// @Summary Get all reservations
// @Description Mendapatkan semua reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations [get]
func (h *reservationHandler) GetReservations(c echo.Context) error {
	userIDinterface := c.Get("user_id")
	userIDfloat, ok := userIDinterface.(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
	}
	userID := int(userIDfloat)
	isAdmin, ok := c.Get("is_admin").(bool)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("Unauthorized"))
	}

	reservations := []*entity.Reservation{}
	var err error
	if !isAdmin {
		reservations, err = h.useCase.GetByUserID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
		}
	} else {
		reservations, err = h.useCase.GetAll(c.Request().Context())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get reservations", reservations))
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
