package middleware

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
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
	for _, r := range roles {
		if hasRole(r, token.Claims.(jwt.MapClaims)["roles"].([]model.Role)) {
			return c.Next()
		}
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Action not allowed", "data": nil})
}

// return true, if a slice of roles contain the requested role tag (as in the jwt claims)
func hasRole(queryRole model.Role, roles []model.Role) bool {
	for _, role := range roles {
		if role.Tag == queryRole.Tag {
			return true
		}
	}
	return false
}
