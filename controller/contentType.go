package controller

import (
	"context"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init ContentTypes: 'blogposts' and 'events'
var (
	blogpost = bson.D{
		{Key: "typename", Value: "blogpost"},
		{Key: "collection", Value: "blogposts"},
		{Key: "field_schema", Value: bson.M{
			"description": "string",
			"text":        "string",
		},
		},
	}

	event = bson.D{
		{Key: "typename", Value: "event"},
		{Key: "collection", Value: "events"},
		{Key: "field_schema", Value: bson.M{
			"description": "string",
			"date":        "time.Time",
			"place":       "string",
		},
		},
	}
)

// Initialize Collection ContentTypes with 'blogposts' and 'events'
func InitContentTypes() error {
	_, err := GetContentType(bson.D{{Key: "typename", Value: "blogpost"}})
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("contenttypes").InsertOne(ctx, blogpost); err != nil {
			return err
		}
	}
	_, err = GetContentType(bson.D{{Key: "typename", Value: "event"}})
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("contenttypes").InsertOne(ctx, event); err != nil {
			return err
		}
	}
	err = nil
	return err
}

// Return all ContentTypes that match the filter
func GetContentTypes(filter interface{}) ([]*model.ContentType, error) {
	// A slice of tasks for storing the decoded documents
	var result []*model.ContentType

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("contenttypes").Find(ctx, filter)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ct model.ContentType
		err := cursor.Decode(&ct)
		if err != nil {
			return result, err
		}

		result = append(result, &ct)
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}

// Return a single ContentType that matches the filter
func GetContentType(filter interface{}) (*model.ContentType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ct *model.ContentType
	err := database.DB.Collection("contenttypes").FindOne(ctx, filter).Decode(&ct)
	if err != nil {
		return nil, err
	}

	return ct, nil
}

// Return a single ContentType with provided ID
func GetContentTypeById(id string) (*model.ContentType, error) {
	ctID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": ctID}
	return GetContentType(filter)
}

// Insert content type with provided Parameters in DB
func CreateContentType(ct *model.ContentType) (*mongo.InsertOneResult, error) {
	// Initialize metadata
	ct.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").InsertOne(ctx, ct)
}

// Delete content type with provided ID in DB
func DeleteContentType(id string) (*mongo.DeleteResult, error) {
	ctID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": ctID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").DeleteOne(ctx, filter)
}

// Validator

func IsValidContentCollection(coll string) bool {
	filter := bson.M{"collection": coll}
	if _, err := GetContentType(filter); err != nil {
		return false
	} else {
		return true
	}
}
