package handler

import (
	"context"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"e-meeting-api/presenter/response"
	"fmt"
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

func (h *userHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse("invalid user id"))
	}

	user, err := h.useCase.GetUser(c.Request().Context(), strconv.Itoa(id))

	if err != nil {
		return c.JSON(http.StatusNotFound, response.NotFoundResponse("user not found"))
	}
	// userResponse := model.UserFromEntity(user)

	return c.JSON(http.StatusOK, response.SuccessResponse("success get user", user))
}

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

	err = h.useCase.UpsertUser(context.Background(), user)

	if err != nil {
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusBadRequest, response.BadRequestResponse("email already exists"))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error", http.StatusInternalServerError))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("register succesfully", user))
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	user := &model.UserRequest{}
	fmt.Println(user)
	err := c.Bind(user)
	if err != nil {
		return c.JSON(500, err)
	}
	err = h.useCase.UpsertUser(context.Background(), user)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, user)
}

func (h *userHandler) Login(c echo.Context) error {
	userReq := &model.UserRequest{}

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

func (h *userHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, response.SuccessResponse("success logout", nil))
}
