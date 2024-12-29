package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

// GetRooms meng-handle request untuk mendapatkan daftar rooms
// @Summary Get all rooms
// @Description Get rooms with filters (name, roomType, capacity)
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "Room name filter"
// @Param roomType query int false "Room type filter (1=Small, 2=Medium, 3=Large)"
// @Param capacity query int false "Capacity filter (1=<10, 2=11-50, 3=51-100)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms [get]
func (h *roomHandler) GetRooms(c echo.Context) error {
	// Ambil parameter query dari request
	name := c.QueryParam("name")
	roomType := c.QueryParam("roomType")
	capacity := c.QueryParam("capacity")

	// Default nilai untuk roomType
	roomTypeInt := 0

	// Jika roomType ada, coba konversi ke integer
	if roomType != "" {
		var err error
		roomTypeInt, err = strconv.Atoi(roomType)
		if err != nil || roomTypeInt < 1 || roomTypeInt > 3 {
			return c.JSON(http.StatusBadRequest, "Invalid roomType")
		}
	}

	// Default nilai untuk capacity
	capacityInt := 0

	// Jika capacity ada, coba konversi ke integer
	if capacity != "" {
		var err error
		capacityInt, err = strconv.Atoi(capacity)
		if err != nil || capacityInt < 1 || capacityInt > 3 {
			// Jika ada error saat konversi atau nilai tidak valid (1-3), kembalikan error
			return c.JSON(http.StatusBadRequest, "Invalid capacity")
		}
	}

	// Panggil use case untuk mendapatkan rooms berdasarkan filter
	rooms, err := h.roomUseCase.GetAll(c.Request().Context(), name, roomTypeInt, capacityInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Return response dengan data rooms
	return c.JSON(http.StatusOK, rooms)
}

// CreateRoom
// @Summary Create room
// @Description Menyimpan room dengan upload file gambar
// @Tags Room
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param name formData string true "Room Name"
// @Param room_type_id formData int true "Room Type ID"
// @Param price formData int true "Price"
// @Param capacity formData int true "Capacity"
// @Param image formData file true "Room Image"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms [post]
func (h *roomHandler) SaveRoom(c echo.Context) error {

	fmt.Println("Received Headers:", c.Request().Header)
	room := &model.RoomRequest{}

	room.Name = c.FormValue("name")
	room.RoomType, _ = strconv.Atoi(c.FormValue("room_type_id"))
	room.Price, _ = strconv.Atoi(c.FormValue("price"))
	room.Capacity, _ = strconv.Atoi(c.FormValue("capacity"))

	fmt.Printf("RoomRequest after manual binding: %+v\n", room)

	err := c.Bind(room)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request body"))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid file"))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("failed to open file"))
	}
	defer src.Close()

	// Validasi file gambar
	isValid, err := util.ValidateImageFile(src)
	if !isValid {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse(err.Error()))
	}
	// Jika file ada, proses upload
	ext := filepath.Ext(file.Filename)
	randomFileName := util.GenerateRandomString(25) + ext
	uploadPath := "public/photos/" + randomFileName

	err = util.SaveUploadedFile(file, uploadPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to upload image"))
	}

	baseURL := os.Getenv("BASE_URL")
	imgUrl := baseURL + "/photos/" + randomFileName
	room.ImgUrl = imgUrl

	err = h.roomUseCase.SaveRoom(c.Request().Context(), room)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
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

	oldRoom, err := h.roomUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to retrieve room details"))
	}

	fmt.Println("Received Headers:", c.Request().Header)
	room := &model.RoomRequest{}

	room.Name = c.FormValue("name")
	room.RoomType, _ = strconv.Atoi(c.FormValue("room_type_id"))
	room.Price, _ = strconv.Atoi(c.FormValue("price"))
	room.Capacity, _ = strconv.Atoi(c.FormValue("capacity"))

	file, err := c.FormFile("image")
	if err == nil {

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("failed to open file"))
		}
		defer src.Close()

		// Validasi file gambar
		isValid, err := util.ValidateImageFile(src)
		if !isValid {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse(err.Error()))
		}

		src.Seek(0, 0)
		// Jika file ada, proses upload
		ext := filepath.Ext(file.Filename)
		randomFileName := util.GenerateRandomString(25) + ext
		uploadPath := "public/photos/" + randomFileName

		err = util.SaveUploadedFile(file, uploadPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to upload image"))
		}

		if oldRoom.ImgUrl != "" {
			oldFilePath := strings.Replace(oldRoom.ImgUrl, os.Getenv("BASE_URL")+"/photos/", "public/photos/", 1)
			if err := os.Remove(oldFilePath); err != nil {
				fmt.Println("Failed to delete old file:", err)
			}
		}

		baseURL := os.Getenv("BASE_URL")
		imgUrl := baseURL + "/photos/" + randomFileName
		room.ImgUrl = imgUrl
	} else {
		// Jika tidak ada file, gunakan gambar lama dari database
		existingRoom, err := h.roomUseCase.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to retrieve existing room details"))
		}
		room.ImgUrl = existingRoom.ImgUrl
	}

	err = h.roomUseCase.UpdateRoom(c.Request().Context(), id, room)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse(err.Error()))
	}

	updatedRoom, err := h.roomUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to retrieve updated room details"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("room updated successfully", updatedRoom))
}

// DeleteRoom
// @Summary Delete room by id
// @Description Menghapus room berdasarkan id
// @Tags Room
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Room ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /rooms/{id} [delete]
func (h *roomHandler) DeleteRoom(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid room id"))
	}
	err = h.roomUseCase.DeleteRoom(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to delete room"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("room deleted successfully", nil))
}
