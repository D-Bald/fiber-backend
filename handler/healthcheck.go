package handler

import (
	"github.com/gofiber/fiber/v2"
)

func Healthcheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Fiber-Backend up and running"})
}
