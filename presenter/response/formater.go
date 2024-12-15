package response

import "net/http"

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, status int) *APIResponse {
	return &APIResponse{
		Status:  status,
		Message: message,
		Data:    nil,
	}
}

func NotFoundResponse(message string) *APIResponse {
	return &APIResponse{
		Status:  http.StatusNotFound,
		Message: message,
		Data:    nil,
	}
}

func BadRequestResponse(message string) *APIResponse {
	return &APIResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Data:    nil,
	}
}

func UnauthorizedResponse(message string) *APIResponse {
	return &APIResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
		Data:    nil,
	}
}
