package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type roomtypeHandler struct {
	roomtypeUseCase usecase.RoomtypeUseCase
}

func NewRoomTypeHandler(roomtypeUseCase usecase.RoomtypeUseCase) RoomtypeHandler {
	return &roomtypeHandler{roomtypeUseCase: roomtypeUseCase}
}

// GetRoomTypes
// @Summary Get all room types
// @Description Mendapatkan semua jenis room
// @Tags Room Type
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /roomtypes [get]
func (h *roomtypeHandler) GetRoomTypes(c echo.Context) error {
	roomTypes, err := h.roomtypeUseCase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to get room types",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get room types", roomTypes))
}

func (h *roomtypeHandler) GetRoomType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room type id"))
	}
	room_type, err := h.roomtypeUseCase.GetByID(c.Request().Context(), id)
	if room_type == nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("room type not found"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success get room type", room_type))
}
