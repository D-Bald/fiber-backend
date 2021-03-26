package handler

import (
	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Content Entries
func GetAllContentEntries(c *fiber.Ctx) error {
	db := database.Mg.Db
	var entries []model.Content
	db.Find(&entries)
	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": entries})
}

// GetContent query content
func GetContent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.Mg.Db
	var content model.Content
	db.Find(&content, id)
	if content.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No content found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": content})
}

// CreateContent new content
func CreateContent(c *fiber.Ctx) error {
	db := database.Mg.Db
	content := new(model.Content)
	if err := c.BodyParser(content); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create content", "data": err})
	}
	db.Create(&content)
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": content})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.Mg.Db

	var content model.Content
	db.First(&content, id)
	if content.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No content found with ID", "data": nil})

	}
	db.Delete(&content)
	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "data": nil})
}
