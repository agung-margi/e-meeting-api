package handler

import (
	"context"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/model"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) UserHandler {
	return &userHandler{useCase: useCase}
}

func (h *userHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.useCase.GetUser(context.Background(), id)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, user)
}

func (h *userHandler) SaveUser(c echo.Context) error {
	user := &model.User{}
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

func (h *userHandler) UpdateUser(c echo.Context) error {
	user := &model.User{}
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
