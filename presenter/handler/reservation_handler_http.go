package handler

import (
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("unauthorized"))
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

	userId := c.Get("user_id").(float64)

	var req model.ReservationRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}

	req.UserID = int(userId)

	reservation, err := h.useCase.SaveReservation(c.Request().Context(), &req, int(userId))
	if err != nil {
		log.Printf("Error saving reservation: %v", err)

		if strings.Contains(err.Error(), "invalid reservation start time") {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse(err.Error()))
		}
		if strings.Contains(err.Error(), "room is not available") {
			return c.JSON(http.StatusConflict, response.ErrorResponse(err.Error(), http.StatusConflict))
		}

		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to save reservation"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("New Reservation Success Added", reservation))
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

// CancelReservation
// @Summary Cancel reservation
// @Description Membatalkan reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations/{id}/cancel [put]
func (h *reservationHandler) CancelReservation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid reservation id"))
	}

	userId := c.Get("user_id").(float64)
	isAdmin := c.Get("is_admin").(bool)

	err = h.useCase.PayReservation(c.Request().Context(), id, int(userId), isAdmin)
	if err != nil {
		switch err.Error() {
		case "unauthorized":
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("unauthorized"))
		case "reservation that has already been paid":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot cancel a reservation that has already been paid"))
		case "reservation that has already started":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot cancel a reservation that has already started"))
		case "reservation that has already ended":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot cancel a reservation that has already ended"))
		case "reservation that has already been cancelled":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot cancel a reservation that has already been cancelled"))
		default:
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to cancel reservation"))
		}
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success cancel reservation", nil))
}

// PayReservation
// @Summary Pay reservation
// @Description Melakukan pembayaran reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations/{id}/pay [put]
func (h *reservationHandler) PayReservation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid reservation id"))
	}

	userId := c.Get("user_id").(float64)
	isAdmin := c.Get("is_admin").(bool)

	err = h.useCase.PayReservation(c.Request().Context(), id, int(userId), isAdmin)
	if err != nil {
		switch err.Error() {
		case "unauthorized":
			return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("unauthorized"))
		case "reservation that has already been paid":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot pay for a reservation that has already been paid"))
		case "reservation that has already started":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot pay for a reservation that has already started"))
		case "reservation that has already ended":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot pay for a reservation that has already ended"))
		case "reservation that has already been cancelled":
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("cannot pay for a reservation that has already been cancelled"))
		default:
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to process payment"))
		}
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success pay reservation", nil))
}
