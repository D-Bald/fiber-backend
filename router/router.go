package router

import (
	"github.com/D-Bald/fiber-backend/handler"
	"github.com/D-Bald/fiber-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// API Route
	api := app.Group("/api", logger.New())

	// Sample Endpoint
	sample := api.Group("/sample")
	sample.Get("/", handler.GetSample)
	sample.Post("/", handler.CreateSample)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	user.Get("/", handler.GetUsers)
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// ContentTypes
	contentTypes := api.Group("/contenttypes") // Insert middleware.Protected() here after testing
	contentTypes.Get("/", handler.GetAllContentTypes)
	contentTypes.Get("/:id", handler.GetContentType)
	contentTypes.Post("/", handler.CreateContentType)
	contentTypes.Delete("/:id", handler.DeleteContentType)

	// Content
	content := api.Group("/:type", func(c *fiber.Ctx) error { //'type' must be a Collection name
		if handler.ValidCollection(c.Params("type")) {
			return c.Next()
		} else {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Review your input for Content Type Collection", "data": nil})
		}
	})
	content.Get("/", handler.GetAllContentEntries)
	content.Get("/:id", handler.GetContent)
	content.Post("/", middleware.Protected(), handler.CreateContent) // Must be changed, if non-users should be able to leave Comments. || 'Comments' API?
	content.Delete("(/:id", middleware.Protected(), handler.DeleteContent)
}
