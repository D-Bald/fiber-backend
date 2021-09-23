package controller

import (
	"context"
	"reflect"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Content struct for Mongo Driver
// Shadowing the ID field of model.Content with ObjectID instead of string type
type ContentMongo struct {
	Content *model.Content
	ID      primitive.ObjectID `bson:"_id" json:"_id" xml:"_id" form:"_id"`
}

// Return the Content Struct and change ID to string
func (cm *ContentMongo) toContent() *model.Content {
	content := cm.Content
	content.ID = cm.ID.Hex()

	return content
}

// Return all content entries from collection coll that match the filter
func GetContent(coll string, filter interface{}) ([]*model.Content, error) {
	var result []*model.Content

	// if filter contains id string, parses it to ObjectID
	v := reflect.ValueOf(filter)
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			// Differentiates between different fields of the struct specified by their bson flag
			if v.Type().Field(i).Tag.Get("bson") == "_id" {
				contentEntry, err := GetContentById(coll, v.Field(i).String())
				if err != nil {
					return nil, err
				}
				result = append(result, contentEntry)
				return result, nil
				// If "_id" is in the query input, the database result ist handed back here, and wie don't have to deal with other query fields. If it is not in the query input, we don't have to deal with parsing the id to an ObjectID in the following lines.
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection(coll).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cMongo ContentMongo
		err := cursor.Decode(&cMongo)
		if err != nil {
			return nil, err
		}
		content := cMongo.toContent()
		result = append(result, content)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}

// Return a single content entry from collection coll that matches the filter. Filter must be structured in bson types.
func GetContentEntry(coll string, filter interface{}) (*model.Content, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cMongo *ContentMongo
	err := database.DB.Collection(coll).FindOne(ctx, filter).Decode(&cMongo)
	if err != nil {
		return nil, err
	}
	result := cMongo.toContent()
	return result, nil
}

// Return Content from collection coll with provided ID
func GetContentById(coll string, id string) (*model.Content, error) {
	cID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": cID}
	return GetContentEntry(coll, filter)
}

// Insert content entry in collection coll with provided Parameters
func CreateContent(coll string, content *model.Content) (*mongo.InsertOneResult, error) {
	// Get corresponding content type set the ContentTypeID reference.
	// ct's FieldSchema could be accessed for validation
	ct, err := GetContentType(bson.M{"collection": coll})
	if err != nil {
		return new(mongo.InsertOneResult), err
	}

	// Initialize metadata
	id := primitive.NewObjectID()
	content.Init(*ct)

	// Create and insert ContenMongo instance
	mongoContent := ContentMongo{Content: content, ID: id}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection(coll).InsertOne(ctx, mongoContent)
}

// Update content entry in collection coll with provided parameters
func UpdateContent(coll string, id string, input *model.ContentUpdate) (*mongo.UpdateResult, error) {
	cID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return new(mongo.UpdateResult), err
	}
	// Update content with provided ID and sets field value `updatet_at`
	filter := bson.M{"_id": cID}
	update := bson.D{
		{Key: "$set", Value: *input},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection(coll).UpdateOne(ctx, filter, update)
}

// Delete content entry provided ID in DB
func DeleteContent(coll string, id string) (*mongo.DeleteResult, error) {
	cID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": cID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection(coll).DeleteOne(ctx, filter)
}
