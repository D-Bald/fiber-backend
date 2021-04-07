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
	user.Post("/", handler.CreateUser)
	user.Get("/:id", handler.GetUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// ContentTypes
	contentTypes := api.Group("/contenttypes")
	contentTypes.Get("/", handler.GetAllContentTypes)
	contentTypes.Post("/", middleware.Protected(), handler.CreateContentType)
	contentTypes.Get("/:id", handler.GetContentType)
	contentTypes.Delete("/:id", middleware.Protected(), handler.DeleteContentType)

	// Content
	content := api.Group("/:content", func(c *fiber.Ctx) error { // `content` has to be a collection
		if controller.IsValidContentCollection(c.Params("content")) {
			return c.Next()
		} else {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Review your route for valid content type", "data": nil})
		}
	})
	content.Get("/", handler.GetAllContentEntries)
	content.Post("/", middleware.Protected(), handler.CreateContent) // Protection must be changed, if non-users should be able to leave comments => Add admin only middleware here.
	// Query contents by different Paramters
	content.Get("/*", handler.GetContent)
	// content.Get("/:id", handler.GetContentById) // Deprecated: content.Get("/id=:id") is used instead

	content.Patch("/:id", handler.UpdateContent) // Add middleware.Protected() after testing. => Add admin only middleware here
	content.Delete("/:id", middleware.Protected(), handler.DeleteContent)
}
