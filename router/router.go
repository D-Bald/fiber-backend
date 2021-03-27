package router

import (
	"github.com/D-Bald/fiber-backend/handler"
	"github.com/D-Bald/fiber-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

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

	// Content
	content := api.Group("/content")
	content.Get("/", handler.GetAllContentEntries)
	content.Get("/:id", handler.GetContent)
	content.Post("/", middleware.Protected(), handler.CreateContent)
	content.Delete("/:id", middleware.Protected(), handler.DeleteContent)
}
