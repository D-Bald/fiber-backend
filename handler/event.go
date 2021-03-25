package handler

import (
	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Event Entries
func GetAll(c *fiber.Ctx) error {
	db := database.DB
	var entries []model.Event
	db.Find(&entries)
	return c.JSON(fiber.Map{"status": "success", "message": "All Event Entries", "data": entries})
}

// GetEvent query event
func GetEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var event model.Event
	db.Find(&event, id)
	if event.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No event found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Event found", "data": event})
}

// CreateEvent new event
func CreateEvent(c *fiber.Ctx) error {
	db := database.DB
	event := new(model.Event)
	if err := c.BodyParser(event); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create event", "data": err})
	}
	db.Create(&event)
	return c.JSON(fiber.Map{"status": "success", "message": "Created event", "data": event})
}

// DeleteEvent delete event
func DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var event model.Event
	db.First(&event, id)
	if event.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No event found with ID", "data": nil})

	}
	db.Delete(&event)
	return c.JSON(fiber.Map{"status": "success", "message": "Event successfully deleted", "data": nil})
}
