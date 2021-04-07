package handler

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Content Entries
func GetAllContentTypes(c *fiber.Ctx) error {
	result, err := controller.GetContentTypes(bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Content Types", "data": result})
}

// GetContent query content
func GetContentType(c *fiber.Ctx) error {
	ct, err := controller.GetContentTypeById(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "data": err.Error()})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content Type found", "data": ct})
}

// CreateContent
func CreateContentType(c *fiber.Ctx) error {
	type NewContentType struct {
		ID          primitive.ObjectID     `bson:"_id" json:"_id"`
		TypeName    string                 `bson:"typename" json:"typename"`
		Collection  string                 `bson:"collection" json:"collection"`
		FieldSchema map[string]interface{} `bson:"field_schema" json:"field_schema"`
	}

	ct := new(model.ContentType)
	// Parse input
	if err := c.BodyParser(ct); err != nil || ct.TypeName == "" || ct.Collection == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Review your input: 'typename' and 'collection' required", "data": err.Error()})
	}

	// Check if already exists
	checkTypeName, _ := controller.GetContentType(bson.M{"typename": ct.TypeName})
	checkCollection, _ := controller.GetContentType(bson.M{"collection": ct.Collection})
	if checkTypeName != nil || checkCollection != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "data": nil})
	}

	// Insert in DB
	if _, err := controller.CreateContentType(ct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content Type", "data": err.Error()})
	}

	// Response
	newCt := NewContentType{
		ID:          ct.ID,
		TypeName:    ct.TypeName,
		Collection:  ct.Collection,
		FieldSchema: ct.FieldSchema,
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created Content Type", "data": newCt})
}

// DeleteContent delete content
func DeleteContentType(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if content type with given id exists
	if _, err := controller.GetContentTypeById(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "data": err.Error()})
	}

	// Delete in DB
	result, err := controller.DeleteContentType(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content Type", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content Type successfully deleted", "data": result})
}
