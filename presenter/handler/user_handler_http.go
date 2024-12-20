package handler

import (
	"context"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"net/http"
	"strconv"

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
// @Summary Register user
// @Description Menyimpan User
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.UserRequest true "User Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /users [post]
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
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param user body model.UpdateUserRequest true "Update User Request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /users/{id} [put]
func (h *userHandler) UpdateUser(c echo.Context) error {
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

	user := &model.UpdateUserRequest{}

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

	err = h.useCase.UpdateUser(context.Background(), id, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}

	updatedUser, err := h.useCase.GetUser(context.Background(), id) // Fetch the updated user from DB
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
