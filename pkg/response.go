package pkg

import (
	"github.com/go-playground/validator/v10"
	"log"
)

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type successResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SwaggerErrorResponse errorResponse
type SwaggerSuccessResponse successResponse

func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FormatValidationErrors(err error) string {
	var errorMsg string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errorMsg += err.Field() + " is required. "
		case "min":
			errorMsg += err.Field() + " must be at least " + err.Param() + " characters long. "
		case "max":
			errorMsg += err.Field() + " must be no more than " + err.Param() + " characters long. "
		case "gte":
			errorMsg += err.Field() + " must be greater than or equal to " + err.Param() + ". "
		case "gt":
			errorMsg += err.Field() + " must be greater than " + err.Param() + ". "
		default:
			errorMsg += err.Field() + " has an invalid value. "
		}
	}
	return errorMsg
}

func NewErrorResponse(status int, message, details string) errorResponse {
	if status == 0 {
		status = 500
	}

	if message == "" {
		message = "Internal Server Error"
	}
	return errorResponse{
		Status:  status,
		Message: message,
		Details: details,
	}
}

func NewSuccessResponse(status int, message string, data interface{}) successResponse {
	if status == 0 {
		status = 200
	}

	if message == "" {
		message = "Success"
	}

	return successResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
