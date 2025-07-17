package utils

import (
	"codematic/internal/shared/model"

	"github.com/gofiber/fiber/v2"
)

func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SendCreated(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(model.SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SendError(c *fiber.Ctx, statusCode int, message string, errs []model.ErrorDetail) error {
	return c.Status(statusCode).JSON(model.ErrorResponse{
		Status:  "error",
		Message: message,
		Errors:  errs,
	})
}
