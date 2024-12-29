package handler

import (
	"context"
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

type userHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) UserHandler {
	return &userHandler{useCase: useCase}
}

// @host      localhost:8080
// @BasePath  /

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

// GetUser
// @Summary Get user by id
// @Description Mendapatkan User berdasarkan id
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{id} [get]
func (h *userHandler) GetUser(c echo.Context) error {
	// Ambil user_id dari token (context middleware)
	tokenUserID, ok := c.Get("user_id").(float64) // JWT claims biasanya float64
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("invalid token claims"))
	}

	// Konversi user_id dari token ke int
	tokenUserIDInt := int(tokenUserID)

	// Ambil id dari parameter URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid user id"))
	}

	// Validasi: hanya izinkan akses jika token user_id sama dengan id
	if tokenUserIDInt != id {
		return c.JSON(http.StatusForbidden, response.ForbiddenResponse("invalid access"))
	}

	// Ambil data user dari database
	user, err := h.useCase.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("user not found"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get user", user))
}

// SaveUser
// @Summary Register
// @Description Menyimpan User
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.UserRequest true "User Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /register [post]
func (h *userHandler) SaveUser(c echo.Context) error {
	user := &model.UserRequest{}
	err := c.Bind(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request data"))
	}

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid data: "+err.Error()))
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = h.useCase.SaveUser(context.Background(), user)

	if err != nil {
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("email already exists"))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("register succesfully", user))
}

// UpdateUser
// @Summary Update user by id
// @Description Memperbarui User berdasarkan id
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param username formData string false "Username"
// @Param email formData string false "Email"
// @Param password formData string false "Password"
// @Param isAdmin formData string false "IsAdmin"
// @Param isActive formData string false "IsActive"
// @Param image formData file false "User Image"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /users/{id} [put]
func (h *userHandler) UpdateUser(c echo.Context) error {
	// Ambil user_id dari token (context middleware)
	tokenUserID, ok := c.Get("user_id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("invalid token claims"))
	}

	// Konversi user_id dari token ke int
	tokenUserIDInt := int(tokenUserID)

	// Ambil id dari parameter URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid user id"))
	}

	// Validasi: hanya izinkan akses jika token user_id sama dengan id
	if tokenUserIDInt != id {
		return c.JSON(http.StatusForbidden, response.ForbiddenResponse("invalid access"))
	}
	oldUser, err := h.useCase.GetUser(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("user not found"))
	}

	user := &model.UpdateUserRequest{}

	user.Username = c.FormValue("username")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.ImgUrl = c.FormValue("imgUrl")

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request data"))
	}
	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid data: "+err.Error()))
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

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

		if oldUser.ImgUrl != "" {
			oldFilePath := strings.Replace(oldUser.ImgUrl, os.Getenv("BASE_URL")+"/photos/", "public/photos/", 1)
			if err := os.Remove(oldFilePath); err != nil {
				fmt.Println("Failed to delete old file:", err)
			}
		}

		baseURL := os.Getenv("BASE_URL")
		imgUrl := baseURL + "/photos/" + randomFileName
		user.ImgUrl = imgUrl
	} else {
		// Jika tidak ada file, gunakan gambar lama dari database
		existingUser, err := h.useCase.GetUser(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to retrieve existing room details"))
		}
		user.ImgUrl = existingUser.ImgUrl
	}

	err = h.useCase.UpdateUser(context.Background(), id, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	updatedUser, err := h.useCase.GetUser(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("error fetching user", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success update user", updatedUser))
}

// Login
// @Summary Login user
// @Description Melakukan login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "Login Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /login [post]
func (h *userHandler) Login(c echo.Context) error {
	userReq := &model.LoginRequest{}

	if err := c.Bind(userReq); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request data"))
	}

	token, _, err := h.useCase.Login(context.Background(), userReq.Username, userReq.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.UnauthorizedResponse("username or password invalid"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success login", map[string]interface{}{
		"token": token,
	}))
}

// Logout
// @Summary Logout user
// @Description Melakukan logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /logout [post]
func (h *userHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, response.SuccessResponse("success logout", nil))
}
