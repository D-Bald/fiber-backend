package middleware

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/form3tech-oss/jwt-go"

	"github.com/gofiber/fiber/v2"
)

// This Middleware requires Protected() to be called in middleware chain before.
// Rolechecker checks "roles" claim in jwt token against permissions of the contenttype of the requested content
// Checks also for "admin" claim and passes if it is true
func ApplyPermissions(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["admin"].(bool) {
		return c.Next()
	}
	ct, _ := controller.GetContentTypeByCollection(c.Params("content")) // Error check obsolet, because IsValidContentCollection is called before.
	roles := ct.Permissions[c.Method()]
	for _, rID := range roles {
		if hasRole(rID.Hex(), token.Claims.(jwt.MapClaims)["roles"].([]interface{})) {
			return c.Next()
		}
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Action not allowed", "data": nil})
}

// return true, if a slice of roles contain the requested role ID (as in the jwt claims)
func hasRole(rID string, roles []interface{}) bool {
	for _, elem := range roles {
		if elem == rID {
			return true
		}
	}
	return false
}
