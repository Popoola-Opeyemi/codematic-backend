package utils

import (
	"codematic/internal/shared/model"

	"github.com/gofiber/fiber/v2"
)

func ExtractUserIDFromJWT(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*model.Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.UserID
}

func ExtractTenantFromJWT(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*model.Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.TenantID
}
