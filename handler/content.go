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
	var ct *model.Content

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DB.Collection(coll).FindOne(ctx, filter).Decode(&ct)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func GetAllContentEntries(c *fiber.Ctx) error {
	var result []*model.Content
	coll := c.Params("content")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection(coll).Find(ctx, bson.D{})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content found", "data": result})
	}
	defer cursor.Close(ctx)

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

	if len(result) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content found.", "data": result})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": result})
}

// GetContent query content
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	contentID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on ID", "data": err.Error()})
	}
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
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create content", "data": err.Error()})
	}

	// Get corresponding content type set the ContentType reference.
	// ct's FieldSchema could be accessed for validation
	col := c.Params("content")
	ct, err := getContentType(bson.D{{Key: "collection", Value: col}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create content", "data": err.Error()})
	}

	// Initialise metadata
	content.ID = primitive.NewObjectID()
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()
	content.ContentType = ct.ID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := database.DB.Collection(col).InsertOne(ctx, &content); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": content})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	contentID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on ID", "data": err.Error()})
	}
	filter := bson.D{{Key: "_id", Value: contentID}}
	if _, err := getContent(coll, filter); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": nil})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.DB.Collection(coll).DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "data": result})
}
