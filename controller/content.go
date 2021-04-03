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

// Return all content entries from collection coll that match the filter
func GetContentEntries(coll string, filter interface{}) ([]*model.Content, error) {
	var result []*model.Content

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection(coll).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var con model.Content
		err := cursor.Decode(&con)
		if err != nil {
			return nil, err
		}

		result = append(result, &con)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}

// Return a single content entry from collection coll that matches the filter
func GetContent(coll string, filter interface{}) (*model.Content, error) {
	var ct *model.Content

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DB.Collection(coll).FindOne(ctx, filter).Decode(&ct)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

// Return Content from collection coll with given ID
func GetContentById(coll string, id string) (*model.Content, error) {
	cID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: cID}}
	return GetContent(coll, filter)
}

// Insert content entry in collection coll with given Parameters
func CreateContent(coll string, content *model.Content) (*mongo.InsertOneResult, error) {
	// Get corresponding content type set the ContentType reference.
	// ct's FieldSchema could be accessed for validation
	ct, err := GetContentType(bson.D{{Key: "collection", Value: coll}})
	if err != nil {
		return new(mongo.InsertOneResult), err
	}

	// Initialize metadata
	content.Init(ct.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection(coll).InsertOne(ctx, content)
}

// Update content entry in collection coll with given Parameters
// ADD CONTROLLER FOR PATCH HANDLER HERE

// Delete content entry given ID in DB
func DeleteContent(coll string, id string) (*mongo.DeleteResult, error) {
	cID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: cID}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection(coll).DeleteOne(ctx, filter)
}
