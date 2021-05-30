package main

import (
	"fmt"
	"log"

	"github.com/D-Bald/fiber-backend/config"
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

	// Initialize Role System
	if err := controller.InitRoles(); err != nil {
		log.Fatal(err)
	}

	// Initialize content types
	if err := controller.InitContentTypes(); err != nil {
		log.Fatal(err)
	}

	// Initialize admin user
	if err := controller.InitAdminUser(); err != nil {
		log.Fatal(err)
	}

	// Start app
	router.SetupRoutes(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.Config("FIBER_PORT"))))
}
