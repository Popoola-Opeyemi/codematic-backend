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

// ExtractUserRoleFromJWT extracts the user role from JWT claims
func ExtractUserRoleFromJWT(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*model.Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.Role
}

// ExtractUserEmailFromJWT extracts the user email from JWT claims
func ExtractUserEmailFromJWT(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*model.Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.Email
}

// HasRole checks if the user has the specified role
func HasRole(c *fiber.Ctx, requiredRole model.UserRole) bool {
	userRole := ExtractUserRoleFromJWT(c)
	return userRole == requiredRole.String()
}

// HasAnyRole checks if the user has any of the specified roles
func HasAnyRole(c *fiber.Ctx, requiredRoles ...model.UserRole) bool {
	userRole := ExtractUserRoleFromJWT(c)
	for _, role := range requiredRoles {
		if userRole == role.String() {
			return true
		}
	}
	return false
}

// RequireRole middleware function to check if user has required role
func RequireRole(requiredRole model.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !HasRole(c, requiredRole) {
			return SendErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
		}
		return c.Next()
	}
}

// RequireAnyRole middleware function to check if user has any of the required roles
func RequireAnyRole(requiredRoles ...model.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !HasAnyRole(c, requiredRoles...) {
			return SendErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
		}
		return c.Next()
	}
}
