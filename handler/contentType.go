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

// Initialize Collection ContentTypes with 'blogposts' and 'events'
var (
	blogpost = bson.D{
		{Key: "typename", Value: "blogpost"},
		{Key: "collection", Value: "blogposts"},
		{Key: "field_schema", Value: bson.M{
			"Description": new(string),
			"text":        new(string),
		},
		},
	}

	event = bson.D{
		{Key: "typename", Value: "event"},
		{Key: "collection", Value: "events"},
		{Key: "field_schema", Value: bson.M{
			"Description": new(string),
			"date":        new(time.Time),
		},
		},
	}
)

func InitContentTypes() error {
	// TODO
	// same Code as in CreateContentType but with var "blogpost" and "event"
	err := new(error)
	return *err
}

// return ContentType by given Filter
func getContentType(filter interface{}) (*model.ContentType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ct *model.ContentType
	err := database.DB.Collection("contenttypes").FindOne(ctx, filter).Decode(&ct)
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func IsValidContentCollection(col string) bool {
	filter := bson.D{{Key: "collection", Value: col}}
	if _, err := getContentType(filter); err != nil {
		return false
	} else {
		return true
	}
}

// GetAll query all Content Entries
func GetAllContentTypes(c *fiber.Ctx) error {
	var result []*model.ContentType

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("contenttypes").Find(ctx, bson.D{})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content Types found", "data": result})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var t model.ContentType
		err := cursor.Decode(&t)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": result})
		}

		result = append(result, &t)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": result})
	}

	if len(result) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Content Type found.", "data": result})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Types", "data": result})
}

// GetContent query content
func GetContentType(c *fiber.Ctx) error {
	typeID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on ID", "data": err.Error()})
	}
	filter := bson.D{{Key: "_id", Value: typeID}}
	ct, err := getContentType(filter)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "data": err.Error()})

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
	if err := c.BodyParser(ct); err != nil || ct.TypeName == "" || ct.Collection == "" {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input: 'typename' and 'collection' required", "data": err.Error()})
	}

	checkTypeName, _ := getContentType(bson.D{{Key: "typename", Value: ct.TypeName}})
	checkCollection, _ := getContentType(bson.D{{Key: "collection", Value: ct.Collection}})
	if checkTypeName != nil || checkCollection != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "data": nil})
	}

	// Initialise metadata
	ct.ID = primitive.NewObjectID()
	ct.CreatedAt = time.Now()
	ct.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := database.DB.Collection("contenttypes").InsertOne(ctx, &ct); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create Content Type", "data": err.Error()})
	}

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
	typeID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on ID", "data": err.Error()})
	}
	filter := bson.D{{Key: "_id", Value: typeID}}
	if _, err := getContentType(filter); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "data": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.DB.Collection("contenttypes").DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Content Type", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content Type successfully deleted", "data": result})
}
