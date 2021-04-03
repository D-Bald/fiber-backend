package handler

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

func GetAllContentEntries(c *fiber.Ctx) error {
	coll := c.Params("content")

	result, err := controller.GetContentEntries(coll, bson.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": result})
}

// GetContent query content
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	content, err := controller.GetContentById(coll, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": content})
}

// CreateContent new content
// Collection is created by mongoDB automatically on first insert call
func CreateContent(c *fiber.Ctx) error {
	content := new(model.Content)
	// Parse input
	if err := c.BodyParser(content); err != nil || content.Title == "" || content.Fields == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create content", "data": err.Error()})
	}

	// Get collection from route params
	coll := c.Params("content")

	if _, err := controller.CreateContent(coll, content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": content})
}

// Update content entry with parameters from request body
// ADD PATCH HANDLER HERE

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")
	if _, err := controller.GetContentById(coll, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": nil})
	}

	result, err := controller.DeleteContent(coll, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "data": result})
}
