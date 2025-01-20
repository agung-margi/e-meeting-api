package handler

import (
	"context"
	"e-meeting-api/configs"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
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
// @Success 200 {object} response.APIResponse "success register"
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

		if err.Error() == "username already exists" {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("username already exists"))
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
// @Param isAdmin formData bool false "IsAdmin"
// @Param isActive formData bool false "IsActive"
// @Param language formData string false "Language"
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

	// Ambil user lama dari database
	oldUser, err := h.useCase.GetUser(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("user not found"))
	}

	// Ambil parameter dari form data
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	isAdmin, _ := strconv.ParseBool(c.FormValue("isAdmin"))
	isActive, _ := strconv.ParseBool(c.FormValue("isActive"))
	language := c.FormValue("language")

	// Periksa apakah ada file gambar baru
	var imgUrl string
	file, err := c.FormFile("image")
	if err == nil {
		// Proses file gambar jika ada
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

		// Proses upload file
		ext := filepath.Ext(file.Filename)
		randomFileName := util.GenerateRandomString(25) + ext
		uploadPath := configs.AppConfig.PhotosPath + randomFileName

		err = util.SaveUploadedFile(file, uploadPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to upload image"))
		}

		baseURL := configs.AppConfig.BaseURL
		imgUrl = baseURL + "/photos/" + randomFileName

		// Hapus gambar lama jika ada
		if oldUser.ImgUrl != "" {
			oldFilePath := strings.Replace(oldUser.ImgUrl, baseURL+"/photos/", configs.AppConfig.PhotosPath, 1)
			if err := os.Remove(oldFilePath); err != nil {
				return c.JSON(http.StatusInternalServerError, response.InternalServerErrorResponse("failed to delete old image"))
			}
		}
	} else {
		// Jika tidak ada gambar baru, gunakan gambar lama
		imgUrl = oldUser.ImgUrl
	}

	// Update user di usecase
	updatedUser := &model.UpdateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
		IsAdmin:  isAdmin,
		IsActive: isActive,
		Language: language,
		ImgUrl:   imgUrl,
	}

	err = h.useCase.UpdateUser(context.Background(), id, updatedUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	// Ambil data user yang sudah diperbarui
	updatedUserResponse, err := h.useCase.GetUser(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("error fetching user", http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success update user", updatedUserResponse))
}

// Login
// @Summary Login user
// @Description Melakukan login user berdasarkan username dan password misalnya "username": "user1", "password": "password"
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

// ForgetPassword
// @Summary Forget password user
// @Description Melakukan lupa password user berdasarkan email misalnya "email": "6T9X5@example.com" untuk mendapatkan link reset password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.ForgetPasswordRequest true "ResetPassword Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /forget-password [post]
func (h *userHandler) ForgetPassword(c echo.Context) error {
	req := &model.ForgetPasswordRequest{}

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid request data"))
	}
	err := h.useCase.GetResetPassword(context.Background(), req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success send link reset password", nil))
}

// ResetPassword
// @Summary Reset password user
// @Description Melakukan reset password user berdasarkan email, token, dan password misalnya "email": "6T9X5@example.com", "token": "IniToken!~{sa", "password": "passwordbaru"
// @Tags Authentication
// @Accept multipart/form-data
// @Produce json
// @Param token path string true "Reset Password Token"
// @Param email query string true "User Email"
// @Param password formData string true "Password"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /reset-password/{token} [post]
func (h *userHandler) ResetPassword(c echo.Context) error {

	req := &model.ResetPasswordRequest{
		Email:    c.QueryParam("email"),
		Token:    c.Param("token"),
		Password: c.FormValue("password"),
	}
	// Check if the necessary fields are provided
	if req.Email == "" || req.Token == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("Missing required fields", http.StatusBadRequest))
	}

	// Call the use case to reset the password
	err := h.useCase.ResetPassword(context.Background(), req.Email, req.Token, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("Password updated successfully", nil))
}
