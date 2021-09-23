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

// ContentType struct for Mongo Driver
// Shadowing the ID and Permissions fields of model.ContentType with ObjectID for ID and just the string slice of Role tags. Database Entry need not to be updated on role-update then.
type ContentTypeMongo struct {
	ContentType *model.ContentType
	ID          primitive.ObjectID  `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	Permissions map[string][]string `bson:"permissions" json:"permissions" xml:"permissions" form:"permissions"`
}

// Return the ContentType Struct and change ID to string and Permissions to roles
func (ctm *ContentTypeMongo) toContentType() (*model.ContentType, error) {
	ct := ctm.ContentType
	ct.ID = ctm.ID.Hex()
	for permission, roles := range ctm.Permissions {
		for _, tag := range roles {
			role, err := GetRoleByTag(tag)
			if err != nil {
				return nil, err
			}
			ct.Permissions[permission] = append(ct.Permissions[permission], *role)
		}
	}
	return ct, nil
}

// Initialize collection ContentTypes with 'blogposts' and 'events'
func InitContentTypes() error {
	// Checks if blogpost exists
	_, err := GetContentType(bson.D{{Key: "typename", Value: "blogpost"}})
	if err != nil && err == mongo.ErrNoDocuments {
		// Get Roles
		var roles []string
		if defaultRole, err := GetRoleByTag("default"); err != nil {
			return err
		} else {
			roles = append(roles, defaultRole.Tag)
		}
		if adminRole, err := GetRoleByTag("admin"); err != nil {
			return err
		} else {
			roles = append(roles, adminRole.Tag)
		}
		// define blogpost document
		blogpost := bson.D{
			{Key: "typename", Value: "blogpost"},
			{Key: "collection", Value: "blogposts"},
			{Key: "permissions", Value: bson.M{
				"GET":    roles,
				"POST":   roles,
				"PATCH":  roles,
				"DELETE": roles,
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
		var roles []string
		if defaultRole, err := GetRoleByTag("default"); err != nil {
			return err
		} else {
			roles = append(roles, defaultRole.Tag)
		}
		if adminRole, err := GetRoleByTag("admin"); err != nil {
			return err
		} else {
			roles = append(roles, adminRole.Tag)
		}
		// define blogpost document
		event := bson.D{
			{Key: "typename", Value: "event"},
			{Key: "collection", Value: "events"},
			{Key: "permissions", Value: bson.M{
				"GET":    roles,
				"POST":   roles,
				"PATCH":  roles,
				"DELETE": roles,
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

	// Convert empty filter to empty bson.object
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := database.DB.Collection("contenttypes").Find(ctx, filter)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ctMongo ContentTypeMongo
		err := cursor.Decode(&ctMongo)
		if err != nil {
			return result, err
		}
		ct, err := ctMongo.toContentType()
		if err != nil {
			return result, err
		}
		result = append(result, ct)
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

	var ctMongo *ContentTypeMongo
	err := database.DB.Collection("contenttypes").FindOne(ctx, filter).Decode(&ctMongo)
	if err != nil {
		return nil, err
	}
	result, err := ctMongo.toContentType()
	if err != nil {
		return nil, err
	}
	return result, nil
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

// Return a single ContentType with provided typename
func GetContentTypeByTypeName(typename string) (*model.ContentType, error) {
	filter := bson.M{"typename": typename}
	return GetContentType(filter)
}

// Returns a content type basing on a collection
func GetContentTypeByCollection(coll string) (*model.ContentType, error) {
	filter := bson.M{"collection": coll}
	return GetContentType(filter)
}

// Insert content type with provided Parameters in DB
func CreateContentType(ct *model.ContentType) (*mongo.InsertOneResult, error) {
	// Initialize metadata
	id := primitive.NewObjectID()
	ct.Init()

	// Only store role tags
	permissions := make(map[string][]string)
	for perm, roles := range ct.Permissions {
		permissions[perm] = GetRoleTagsFromRoleSlice(roles)
	}

	// Create and insert ContenTypeMongo instance
	mongoCT := ContentTypeMongo{ContentType: ct, ID: id, Permissions: permissions}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").InsertOne(ctx, mongoCT)
}

// Update content type with provided parameters
func UpdateContentType(id string, input *model.ContentTypeInput) (*mongo.UpdateResult, error) {
	// Struct similar to `ContentTypeUpdate` but with ObjectIDs of roles instead of string role names in Permissions
	type ContentTypeUpdateMongo struct {
		ContentTypeUpdate *model.ContentTypeInput
		Permissions       map[string][]string `bson:"permissions,omitempty"` // Shadowing the Permissions field of model.ContentTypeUpdate
	}

	// create Object with ObjectIDs as Roles
	ctUpdate := ContentTypeUpdateMongo{
		ContentTypeUpdate: input,
		Permissions:       make(map[string][]string),
	}

	// Only store role tags
	for perm, roles := range input.Permissions {
		ctUpdate.Permissions[perm] = GetRoleTagsFromRoleSlice(roles)
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
func DeleteRoleFromPermissions(delRole string, ct *model.ContentType) (*mongo.UpdateResult, error) {

	newPermissions := make(map[string][]string)
	for permission, roles := range ct.Permissions {
		for _, r := range roles {
			// all roles of ct except the one to delete are added to the update to keep them
			if r.Tag != delRole {
				newPermissions[permission] = append(newPermissions[permission], r.Tag)
			}
		}
	}
	// Update content type with provided ID and sets field value for `updatet_at`
	ctObjectID, err := primitive.ObjectIDFromHex(ct.ID)
	if err != nil {
		return new(mongo.UpdateResult), err
	}
	filter := bson.M{"_id": ctObjectID}
	update := bson.D{
		{Key: "$set", Value: bson.M{"permissions": newPermissions}},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("contenttypes").UpdateOne(ctx, filter, update)
}

// Returns true if the a contenttype with exists, where the `collection` field value is `coll`
func IsValidContentCollection(coll string) bool {
	filter := bson.M{"collection": coll}
	if _, err := GetContentType(filter); err != nil {
		return false
	} else {
		return true
	}
}

// Returns the Custom fields of a contenttype as map
// Takes the collection of a content type as input
func GetCustomFields(coll string) (map[string]interface{}, error) {
	filter := bson.M{"collection": coll}
	if ct, err := GetContentType(filter); ct != nil {
		return ct.FieldSchema, nil
	} else {
		return nil, err
	}
}
