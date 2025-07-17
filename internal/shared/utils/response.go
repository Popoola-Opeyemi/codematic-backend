package utils

import (
	"codematic/internal/shared/model"

	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
	response := model.SuccessResponse{
		Status: "success",
		Data:   data,
	}

	return c.Status(statusCode).JSON(response)
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	response := model.ErrorResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}

	return c.Status(statusCode).JSON(response)
}

// ResponseCapture captures Fiber response body and status code for middleware use
// Usage: call rc.Capture(c) before c.Next(), then use rc.Body and rc.StatusCode after
// Note: Fiber does not support replacing the response writer, so we capture after c.Next()
type ResponseCapture struct {
	Body       []byte
	StatusCode int
}

func (rc *ResponseCapture) Capture(c *fiber.Ctx) {
	rc.Body = c.Response().Body()
	rc.StatusCode = c.Response().StatusCode()
}
