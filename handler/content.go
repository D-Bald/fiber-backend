package handler

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"github.com/D-Bald/fiber-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Query content entries with filter provided in query params
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	parseObject := new(model.Content)

	// parse input
	if err := c.QueryParser(parseObject); err != nil && err.Error() != "schema: converter not found for primitive.ObjectID" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
	}
	// parse ID manually because fiber's QueryParser has no converter for this type.
	if id := string(c.Request().URI().QueryArgs().Peek("id")); id != "" {
		cID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid ID", "data": err.Error()})
		}
		parseObject.ID = cID
	}

	// Parse custom fields manually
	fields, err := controller.GetCustomFields(coll) // Initialize fields map to avoid nil map error
	parseObject.Fields = make(map[string]interface{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Custom fields not found", "data": err.Error()})
	}
	for f := range fields {
		if fValue := c.Request().URI().QueryArgs().Peek(f); fValue != nil {
			parseObject.Fields[f] = string(fValue)
		}
	}

	// Make filter where slices, maps and structs are parsed inline
	filter, err := utils.MakeQueryFilterFromStruct(parseObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Custom fields not found", "data": err.Error()})
	}

	// get content from DB
	result, err := controller.GetContent(coll, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "content": result})
}

// CreateContent new content
// Collection is created by mongoDB automatically on first insert call
func CreateContent(c *fiber.Ctx) error {
	content := new(model.Content)
	// Parse input
	if err := c.BodyParser(content); err != nil || content.Title == "" || content.Fields == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create content", "content": err.Error()})
	}

	// Get collection from route params
	coll := c.Params("content")

	if _, err := controller.CreateContent(coll, content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "content": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "content": content})
}

// Update content entry with parameters from request body
// Update user with parameters from request body
func UpdateContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")

	uci := new(model.ContentUpdate)
	if err := c.BodyParser(uci); err != nil || uci == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	result, err := controller.UpdateContent(coll, id, uci)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update content entry", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully updated", "result": result})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")
	if _, err := controller.GetContentById(coll, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "result": err.Error()})
	}

	result, err := controller.DeleteContent(coll, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "result": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "result": result})
}
