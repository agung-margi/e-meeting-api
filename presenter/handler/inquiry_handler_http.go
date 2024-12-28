package handler

import (
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type inquiryHandler struct {
	inquiryUsecase usecase.InquiryUsecase
}

func NewInquiryHandler(inquiryUsecase usecase.InquiryUsecase) InquiryHandler {
	return &inquiryHandler{inquiryUsecase: inquiryUsecase}
}

// SaveInquiry
// @Summary Save Inquiry
// @Description Save Inquiry
// @Tags Reservation
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param reservation body model.ReservationRequest true "Reservation Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reservations/inquiry [post]
func (h *inquiryHandler) SaveInquiry(c echo.Context) error {

	userId := c.Get("user_id").(float64)
	userIdInt := int(userId)

	if userIdInt == 0 {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("unauthorized"))
	}

	// Validasi format tanggal

	var req entity.Reservation

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse(err.Error()))
	}

	// Proceed with the rest of your logic
	reservation := &entity.Reservation{
		UserID:       userIdInt,
		RoomID:       req.RoomID,
		BookingDate:  req.BookingDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Name:         req.Name,
		Phone:        req.Phone,
		Company:      req.Company,
		Participants: req.Participants,
		SnackID:      req.SnackID,
		SnackPrice:   req.SnackPrice,
		TotalPrice:   req.TotalPrice,
		Notes:        req.Notes,
	}

	inquiry, err := h.inquiryUsecase.Save(c.Request().Context(), reservation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("inquiry saved successfully", inquiry))
}
