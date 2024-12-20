package handler

import (
	"e-meeting-api/internal/usecase"
	"e-meeting-api/presenter/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type snackHandler struct {
	snackUsecase usecase.SnackUseCase
}

func NewSnackHandler(snackUsecase usecase.SnackUseCase) SnackHandler {
	return &snackHandler{snackUsecase: snackUsecase}
}

// GetSnacks
// @Summary Get Snacks
// @Description Mendapatkan Snacks
// @Tags Snack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /snacks [get]
func (sh *snackHandler) GetSnacks(c echo.Context) error {
	snack, err := sh.snackUsecase.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.SuccessResponse("success get snacks", snack))
}
