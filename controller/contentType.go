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
		if user, err := GetRoleByName("user"); err != nil {
			return err
		} else {
			roles = append(roles, user.ID)
		}
		if admin, err := GetRoleByName("admin"); err != nil {
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
		if user, err := GetRoleByName("user"); err != nil {
			return err
		} else {
			roles = append(roles, user.ID)
		}
		if admin, err := GetRoleByName("admin"); err != nil {
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

// Update content type with provided parameters
func UpdateContentType(id string, input *model.ContentTypeUpdate) (*mongo.UpdateResult, error) {
	// Struct similar to `ContentTypeUpdate` but with ObjectIDs of roles instead of string role names
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
			rObj, err := GetRoleByName(r)
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
	// Drop corresponding collection
	filter := bson.M{"_id": ctID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ct, err := GetContentType(filter)
	if err != nil {
		return nil, err
	}
	err = database.DB.Collection(ct.Collection).Drop(ctx)
	if err != nil {
		return nil, err
	}
	// Delete content type
	return database.DB.Collection("contenttypes").DeleteOne(ctx, filter)
}

// Delete one role from content type permissions.
func DeleteRoleFromPermissions(rID primitive.ObjectID, ct *model.ContentType) (*mongo.UpdateResult, error) {
	permissions := make(map[string][]primitive.ObjectID)
	for permission, roles := range ct.Permissions {
		for _, r := range roles {
			// all roles of ct except the one to delete are added to the update to keep them
			if r != rID {
				permissions[permission] = append(permissions[permission], r)
			}
		}
	}
	// Update content type with provided ID and sets field value for `updatet_at`
	filter := bson.M{"_id": ct.ID}
	update := bson.D{
		{Key: "$set", Value: bson.M{"permissions": permissions}},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").UpdateOne(ctx, filter, update)
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

// Returns the Custom fields of a contenttype as map
func GetCustomFields(coll string) (map[string]interface{}, error) {
	filter := bson.M{"collection": coll}
	if ct, err := GetContentType(filter); ct != nil {
		return ct.FieldSchema, nil
	} else {
		return nil, err
	}
}
