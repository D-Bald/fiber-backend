package router

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/handler"
	"github.com/D-Bald/fiber-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// API Route
	api := app.Group("/api", logger.New())

	// Healthcheck endpoint
	api.Get("/", handler.Healthcheck)

	// Auth endpoints
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User endpoints
	user := api.Group("/user")
	user.Post("/", handler.CreateUser, handler.Login)
	// Query contents by different Paramters
	user.Get("/", middleware.Protected(), handler.GetUsers)

	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// ContentTypes endpoints
	contentTypes := api.Group("/contenttypes")
	contentTypes.Get("/", handler.GetAllContentTypes)
	contentTypes.Post("/", middleware.Protected(), middleware.AdminOnly, handler.CreateContentType)
	contentTypes.Get("/:id", handler.GetContentType)
	contentTypes.Delete("/:id", middleware.Protected(), middleware.AdminOnly, handler.DeleteContentType)

	// Content endpoints
	content := api.Group("/:content", func(c *fiber.Ctx) error { // `content` has to be a collection
		if controller.IsValidContentCollection(c.Params("content")) {
			return c.Next()
		} else {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Review your route for valid content type", "data": nil})
		}
	})
	content.Post("/", middleware.Protected(), middleware.AdminOnly, handler.CreateContent)

	// Query contents by different Paramters
	content.Get("/", handler.GetContent)

	content.Patch("/:id", middleware.Protected(), middleware.AdminOnly, handler.UpdateContent)
	content.Delete("/:id", middleware.Protected(), middleware.AdminOnly, handler.DeleteContent)
}
