package handler

import (
	"context"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

func getContent(coll string, filter interface{}) (*model.Content, error) {
	ctx := context.TODO()
	var ct *model.Content
	err := database.Mg.Db.Collection(coll).FindOne(ctx, filter).Decode(&ct)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func GetAllContentEntries(c *fiber.Ctx) error {
	var result []*model.Content
	coll := c.Params("type")
	ctx := context.TODO()
	cursor, err := database.Mg.Db.Collection(coll).Find(ctx, bson.D{})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content found", "data": result})
	}
	for cursor.Next(ctx) {
		var con model.Content
		err := cursor.Decode(&con)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": result})
		}

		result = append(result, &con)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": result})
	}

	cursor.Close(ctx)

	if len(result) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content found.", "data": result})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": result})
}

// GetContent query content
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("type")
	contentID, err := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.D{{Key: "_id", Value: contentID}}
	content, err := getContent(coll, filter)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": content})
}

// CreateContent new content
// Collection is created by mongoDB automatically on first insert call
func CreateContent(c *fiber.Ctx) error {
	content := new(model.Content)
	if err := c.BodyParser(content); err != nil || content.Title == "" || content.Fields == nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create content", "data": err.Error()})
	}

	// Initialise metadata
	content.ID = primitive.NewObjectID()
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()

	if _, err := database.Mg.Db.Collection(c.Params("type")).InsertOne(context.TODO(), &content); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": content})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("type")
	contentID, err := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.D{{Key: "_id", Value: contentID}}
	if _, err := getContent(coll, filter); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": nil})
	}

	result, err := database.Mg.Db.Collection(coll).DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "data": result})
}

// not used yet used. Could be used for a GET route, that returns the schema, so that the Client knows, which fields to provide
// also interesting for validation: https://docs.gofiber.io/guide/validation
func getSchema(contentType string) (map[string]interface{}, error) {
	filter := bson.D{{Key: "type_name", Value: contentType}}
	ct, err := getContentType(filter)
	if err != nil {
		return nil, err
	}
	return ct.FieldSchema, nil

}
