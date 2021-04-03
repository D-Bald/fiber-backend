package main

import (
	"log"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Create a Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// prevent the server crash from panics like body-parsing invalid input data
	app.Use(recover.New())

	// Connect to the database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	// Initialize Database
	if err := controller.InitContentTypes(); err != nil {
		log.Fatal(err)
	}

	// Initialize admin user
	if err := controller.InitAdminUser(); err != nil {
		log.Fatal(err)
	}

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
