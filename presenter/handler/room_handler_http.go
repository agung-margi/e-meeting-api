package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type roomHandler struct {
	roomUseCase usecase.RoomUseCase
}

func NewRoomHandler(roomUseCase usecase.RoomUseCase) RoomHandler {
	return &roomHandler{roomUseCase: roomUseCase}
}

// @host      localhost:8080
// @BasePath  /

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

// GetRooms
// @Summary Get all rooms
// @Description Mendapatkan semua room
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms [get]
func (h *roomHandler) GetRooms(c echo.Context) error {
	rooms, err := h.roomUseCase.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success get rooms", rooms))
}

// SaveRoom
// @Summary Save room
// @Description Menyimpan room
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param room body model.RoomRequest true "Room Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms [post]
func (h *roomHandler) SaveRoom(c echo.Context) error {
	room := &model.RoomRequest{}
	if err := c.Bind(room); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}
	err := h.roomUseCase.SaveRoom(c.Request().Context(), room)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to save room"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("room saved successfully", room))
}

// GetRoom
// @Summary Get room by id
// @Description Mendapatkan room berdasarkan id
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Room ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms/{id} [get]
func (h *roomHandler) GetRoom(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}
	room, err := h.roomUseCase.GetByID(c.Request().Context(), id)
	if room == nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("room not found"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success get room", room))
}

// UpdateRoom
// @Summary Update room by id
// @Description Mengupdate room berdasarkan id
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Room ID"
// @Param room body model.RoomRequest true "Room Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms/{id} [put]
func (h *roomHandler) UpdateRoom(c echo.Context) error {
	// Parse room ID from URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}

	roomRequest := &model.RoomRequest{}
	if err := c.Bind(roomRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}

	err = h.roomUseCase.UpdateRoom(c.Request().Context(), id, roomRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to update room"))
	}

	updatedRoom, err := h.roomUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to retrieve updated room details"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("room updated successfully", updatedRoom))
}
