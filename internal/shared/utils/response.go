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
