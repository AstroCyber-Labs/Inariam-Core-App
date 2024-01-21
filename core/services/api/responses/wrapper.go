// Package responses provides utility functions for handling HTTP responses.
package responses

import (
	"github.com/labstack/echo/v4"
)

// Error represents a standard error response structure.
type Error struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// Data represents a generic data response structure.
type Data struct {
	Message string `json:"message"`
}

// Response sends a JSON response with the specified status code and data.
func Response(c echo.Context, statusCode int, data interface{}) error {
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// nolint // context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return c.JSON(statusCode, data)
}

// MessageResponse sends a JSON response with a message and the specified status code.
func MessageResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Data{
		Message: message,
	})
}

// ErrorResponse sends a JSON error response with the specified status code and message.
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Error{
		Code:  statusCode,
		Error: message,
	})
}
