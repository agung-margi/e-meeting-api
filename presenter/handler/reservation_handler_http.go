package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/response"
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

// SaveReservation
// @Summary Save reservation
// @Description Menyimpan reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param inquiry_id path int true "Inquiry ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations/{inquiry_id} [post]
func (h *reservationHandler) SaveReservation(c echo.Context) error {
	// Mendapatkan userId dari context (middleware)
	userId := c.Get("user_id").(float64)

	// Mendapatkan id inquiry dari request parameter
	inquiryIdParam := c.Param("inquiry_id")
	inquiryId, err := strconv.Atoi(inquiryIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid inquiry id"))
	}

	// Memanggil use case untuk menyimpan reservation berdasarkan inquiryId dan userId
	reservation, err := h.useCase.Save(c.Request().Context(), inquiryId, int(userId))
	if err != nil {
		if strings.Contains(err.Error(), "inquiry not found") {
			return c.JSON(http.StatusNotFound, response.ErrorResponse("inquiry not found", http.StatusNotFound))
		}
		if strings.Contains(err.Error(), "room is not available") {
			return c.JSON(http.StatusConflict, response.ErrorResponse(err.Error(), http.StatusConflict))
		}

		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("New Reservation Successfully Added", reservation))
}

// GetReservation by id
// @Summary Get reservation by id
// @Description Mendapatkan reservation
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
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

	if len(*reservation) == 0 {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("reservation not found"))
	}

	reservationUserID := (*reservation)[0].UserID

	if reservationUserID != userID && !isAdmin {
		return c.JSON(http.StatusForbidden, response.ErrorResponse("forbidden", http.StatusForbidden))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get reservation", reservation))
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

	err = h.useCase.CancelReservation(c.Request().Context(), id, int(userId), isAdmin)

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
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
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
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
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
// @Param room_type query int false "Room type"
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

	// // var startDate, endDate *time.Time
	// if startDateStr != "" {
	// 	parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid start_date format"))
	// 	}
	// 	startDate = &parsedStartDate
	// }
	// if endDateStr != "" {
	// 	parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid end_date format"))
	// 	}
	// 	endDate = &parsedEndDate
	// }

	roomTypeInt := 0

	if roomType != "" {
		var err error
		roomTypeInt, err = strconv.Atoi(roomType)
		if err != nil || roomTypeInt < 0 {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room_type"))
		}
	}
	if status != "" && status != "paid" && status != "booked" && status != "cancelled" {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid status"))
	}

	reservations, err := h.useCase.GetAll(c.Request().Context(), startDateStr, endDateStr, roomTypeInt, status, filterUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}
	if reservations == nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("reservations not found"))
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

// GetReservationSchedules
// @Summary Get reservation schedules
// @Description Mendapatkan jadwal reservation room berdasarkan rentang tanggal
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations/schedules [get]
func (h *reservationHandler) GetSchedules(c echo.Context) error {

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid start_date format"))
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid end_date format"))
	}

	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("start_date must be before end_date"))
	}

	schedules, err := h.useCase.GetSchedulesByDateRange(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	currentTime := time.Now().UTC()
	for i := range schedules {
		if currentTime.Before(schedules[i].StartTime) {
			schedules[i].Status = "upcoming"
		} else if currentTime.After(schedules[i].EndTime) {
			schedules[i].Status = "done"
		} else {
			schedules[i].Status = "inprogress"
		}
	}

	return c.JSON(http.StatusOK, response.SuccessResponseWithCount("success get schedules", schedules, len(schedules)))
}

// GetDashboardData
// @Summary Get dashboard data
// @Description Mendapatkan data dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /dashboard [get]
func (h *reservationHandler) GetDashboardData(c echo.Context) error {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	if startDateStr == "" {
		startDateStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02") // 30 hari yang lalu
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02") // Hari ini
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid start_date format"))
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid end_date format"))
	}

	if startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("start_date must be before end_date"))
	}

	dashboardData, err := h.useCase.GetDashboardData(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get dashboard data", dashboardData))
}

// ExportReservations
// @Summary Export reservations
// @Description Export reservations
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
// @Router /reservations/export [get]
func (h *reservationHandler) ExportReservations(c echo.Context) error {

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
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("start_date must be before end_date"))
	}

	roomTypeInt := 0

	if roomType != "" {
		var err error
		roomTypeInt, err = strconv.Atoi(roomType)
		if err != nil || roomTypeInt < 0 {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room_type"))
		}
	}
	if status != "" && status != "paid" && status != "booked" && status != "cancelled" {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid status"))
	}

	excelFile, err := h.useCase.ExportReservations(c.Request().Context(), startDateStr, endDateStr, roomTypeInt, status, filterUserID)
	if err != nil {
		if err.Error() == "reservations not found" {
			return c.JSON(http.StatusNotFound, response.NotFoundResponse("reservations not found"))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	// Set header untuk download file Excel
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", `attachment; filename="reservations.xlsx"`)

	// Tulis file ke response
	return excelFile.Write(c.Response().Writer)
}
