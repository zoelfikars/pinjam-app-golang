package util

import "github.com/labstack/echo/v4"

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c echo.Context, httpStatus int, message string, data interface{}) error {
	return c.JSON(httpStatus, StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func FailResponse(c echo.Context, httpStatus int, message string, data interface{}) error {

	return c.JSON(httpStatus, StandardResponse{
		Status:  "fail",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, httpStatus int, message string) error {

	return c.JSON(httpStatus, StandardResponse{
		Status:  "error",
		Message: message,
	})
}
