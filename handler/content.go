package handler

import (
	"fmt"
	"regexp"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Return all content entries in given collection
func GetAllContentEntries(c *fiber.Ctx) error {
	coll := c.Params("content")

	result, err := controller.GetContentEntries(coll, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": result})
}

// Query content by Param 'id'
func GetContentById(c *fiber.Ctx) error {
	coll := c.Params("content")
	content, err := controller.GetContentById(coll, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": content})
}

// Query content entries with filter provided in query params
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	re := regexp.MustCompile(`[a-z\.\_]+=[a-zA-Z0-9\%]+`)
	filterString := re.FindAllString(c.Params("*"), -1)
	filter := make(map[string]interface{})
	for _, v := range filterString {
		v = regexp.MustCompile(`%20`).ReplaceAllString(v, " ")
		s := regexp.MustCompile(`=`).Split(v, -1)
		if s[0] == `id` {
			cID, err := primitive.ObjectIDFromHex(s[1])
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
			}
			filter["_id"] = cID
		} else {
			filter[s[0]] = s[1]
		}
	}
	fmt.Println(filter)
	result, err := controller.GetContentEntries(coll, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": result})
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
// Update user with parameters from request body
func UpdateContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")

	uci := new(model.UpdateContentInput)
	if err := c.BodyParser(uci); err != nil || uci == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	result, err := controller.UpdateContent(coll, id, uci)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update User", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": result})
}

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
