package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @Router /reservations/{id} [post]
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

	reservationUserID := reservation[0].UserID

	if reservationUserID != userID && !isAdmin {
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
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}
	req.UserID = int(userId)

	reservation, err := h.useCase.SaveReservation(c.Request().Context(), &req, int(userId))
	if err != nil {

		if strings.Contains(err.Error(), "invalid reservation start time") {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse(err.Error()))
		}
		if strings.Contains(err.Error(), "room is not available") {
			return c.JSON(http.StatusConflict, response.ErrorResponse(err.Error(), http.StatusConflict))
		}

		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("New Reservation Success Added", reservation))
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

// GetReservations
// @Summary Get all reservations
// @Description Mendapatkan semua reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Param room_type query string false "Room type"
// @Param status query string false "Status"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations [get]
func (h *reservationHandler) GetReservations(c echo.Context) error {
	// Ambil query parameter
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	roomType := c.QueryParam("room_type")
	status := c.QueryParam("status")

	userID := c.Get("user_id")
	isAdmin := c.Get("is_admin")

	if userID == nil || isAdmin == nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("User not authenticated", http.StatusUnauthorized))
	}

	userIDInt, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("User not authenticated", http.StatusUnauthorized))
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("User not authenticated", http.StatusUnauthorized))
	}

	var filterUserID *int
	if !isAdminBool {
		userIDIntConverted := int(userIDInt)
		filterUserID = &userIDIntConverted
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid start_date format"))
		}
		startDate = &parsedStartDate
	}
	if endDateStr != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid end_date format"))
		}
		endDate = &parsedEndDate
	}
	roomTypeInt := 0

	if roomType != "" {
		var err error
		roomTypeInt, err = strconv.Atoi(roomType)
		if err != nil || roomTypeInt < 0 {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room_type"))
		}
	}
	if status != "" && status != "paid" && status != "booked" && status != "cancel" {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid status"))
	}

	reservations, err := h.useCase.GetAll(c.Request().Context(), startDate, endDate, roomTypeInt, status, filterUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success get reservations", reservations))
}

// GetRoomSchedule
// @Summary Get room schedule
// @Description Mendapatkan jadwal room
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param room_id query int true "Room ID"
// @Param date query string true "Date"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /room-schedule [get]
func (h *reservationHandler) GetRoomSchedule(c echo.Context) error {
	// Validate room_id
	roomIdStr := c.QueryParam("room_id")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}

	// Validate date
	date := c.QueryParam("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		log.Println("Invalid date format:", date)
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid date format, expected YYYY-MM-DD"))
	}
	// Fetch reservations
	reservations, err := h.useCase.GetReservationsByRoomAndDate(c.Request().Context(), roomId, date)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseWithCount("success get room schedule", reservations, len(reservations)))
}
