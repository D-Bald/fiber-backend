package middleware

import (
	"github.com/form3tech-oss/jwt-go"

	"github.com/gofiber/fiber/v2"
)

// AdminOnly checks for "admin" claim in jwt token.
// This Middleware requires Protected() to be called in middleware chain before.
func AdminOnly(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["admin"].(bool) {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Admin only", "data": nil})
}
