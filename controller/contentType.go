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

// Initialize collection ContentTypes with 'blogposts' and 'events'
func InitContentTypes() error {
	// Checks if blogpost exists
	_, err := GetContentType(bson.D{{Key: "typename", Value: "blogpost"}})
	if err != nil && err == mongo.ErrNoDocuments {
		// Get Roles
		var roles []primitive.ObjectID
		if user, err := GetRole("user"); err != nil {
			return err
		} else {
			roles = append(roles, user.ID)
		}
		if admin, err := GetRole("admin"); err != nil {
			return err
		} else {
			roles = append(roles, admin.ID)
		}
		// define blogpost document
		blogpost := bson.D{
			{Key: "typename", Value: "blogpost"},
			{Key: "collection", Value: "blogposts"},
			{Key: "permissions", Value: bson.M{
				"get":    roles,
				"post":   roles,
				"patch":  roles,
				"delete": roles,
			}},
			{Key: "field_schema", Value: bson.M{
				"description": "string",
				"text":        "string",
			}},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("contenttypes").InsertOne(ctx, blogpost); err != nil {
			return err
		}
	}
	// checks if event exists
	_, err = GetContentType(bson.D{{Key: "typename", Value: "event"}})
	if err != nil && err == mongo.ErrNoDocuments {
		// Get Roles
		var roles []primitive.ObjectID
		if user, err := GetRole("user"); err != nil {
			return err
		} else {
			roles = append(roles, user.ID)
		}
		if admin, err := GetRole("admin"); err != nil {
			return err
		} else {
			roles = append(roles, admin.ID)
		}
		// define blogpost document
		event := bson.D{
			{Key: "typename", Value: "event"},
			{Key: "collection", Value: "events"},
			{Key: "permissions", Value: bson.M{
				"get":    roles,
				"post":   roles,
				"patch":  roles,
				"delete": roles,
			}},
			{Key: "field_schema", Value: bson.M{
				"description": "string",
				"date":        "time.Time",
				"place":       "string",
			},
			},
		}

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
func GetContentTypes(filter interface{}) ([]*model.ContentTypeOutput, error) {
	// A slice of tasks for storing the decoded documents
	var result []*model.ContentTypeOutput

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

		ctOutput := model.ContentTypeOutput{
			ID:          ct.ID,
			TypeName:    ct.TypeName,
			Collection:  ct.Collection,
			Permissions: make(map[string][]string),
			FieldSchema: ct.FieldSchema,
		}
		// Parse role ObjectIDs to role name strings
		for key, val := range ct.Permissions {
			roles, err := GetRoleNames(val)
			if err != nil {
				return nil, err
			}
			ctOutput.Permissions[key] = roles
		}

		result = append(result, &ctOutput)
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

// Update content type with provided parameters
func UpdateContentType(id string, input *model.ContentTypeUpdate) (*mongo.UpdateResult, error) {
	type mongoContentTypeUpdate struct {
		TypeName    string                          `bson:"typename,omitempty"`
		Collection  string                          `bson:"collection,omitempty"`
		Permissions map[string][]primitive.ObjectID `bson:"permissions,omitempty"`
		FieldSchema map[string]interface{}          `bson:"field_schema,omitempty"`
	}

	// create Object with ObjectIDs as Roles
	ctUpdate := mongoContentTypeUpdate{
		TypeName:    input.TypeName,
		Collection:  input.Collection,
		Permissions: make(map[string][]primitive.ObjectID),
		FieldSchema: input.FieldSchema,
	}

	// Parse role name strings in Permissionsto role ObjectIDs
	for key, val := range input.Permissions {
		var roleObjectIDs []primitive.ObjectID
		for _, r := range val {
			rObj, err := GetRole(r)
			if err != nil {
				return new(mongo.UpdateResult), err
			}
			roleObjectIDs = append(roleObjectIDs, rObj.ID)
		}
		ctUpdate.Permissions[key] = roleObjectIDs
	}

	ctID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return new(mongo.UpdateResult), err
	}
	// Update content type with provided ID and sets field value for `updatet_at`
	filter := bson.M{"_id": ctID}
	update := bson.D{
		{Key: "$set", Value: ctUpdate},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").UpdateOne(ctx, filter, update)
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

// Return true if the a contenttype with exists, where the `collection` field value is `coll`
func IsValidContentCollection(coll string) bool {
	filter := bson.M{"collection": coll}
	if _, err := GetContentType(filter); err != nil {
		return false
	} else {
		return true
	}
}
