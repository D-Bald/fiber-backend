package main

import (
	"log"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/handler"
	"github.com/D-Bald/fiber-backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Create a Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// Connect to the database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	// Initialize Database
	handler.InitContentTypes()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
