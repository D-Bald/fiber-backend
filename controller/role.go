package controller

import (
	"context"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init roles: 'admin' and 'user'
var (
	defaultRole = bson.D{
		{Key: "tag", Value: "default"},
		{Key: "name", Value: "User"},
	}
	admin = bson.D{
		{Key: "tag", Value: "admin"},
		{Key: "name", Value: "Administrator"},
	}
)

// Initialize collection roles with 'admin' and 'user'
func InitRoles() error {
	_, err := GetRoleByTag("default")
	if err != nil && err == mongo.ErrNoDocuments {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("roles").InsertOne(ctx, defaultRole); err != nil {
			return err
		}
	}
	_, err = GetRoleByTag("admin")
	if err != nil && err == mongo.ErrNoDocuments {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("roles").InsertOne(ctx, admin); err != nil {
			return err
		}
	}
	err = nil
	return err
}

// Return all roles that match the filter
func GetRoles(filter interface{}) ([]*model.Role, error) {
	// A slice of tasks for storing the decoded documents
	var result []*model.Role

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert empty filter to empty bson.object
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := database.DB.Collection("roles").Find(ctx, filter)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r model.Role
		err := cursor.Decode(&r)
		if err != nil {
			return result, err
		}

		result = append(result, &r)
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}

// Returns first role that matches the filter
func GetRole(filter interface{}) (*model.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var r *model.Role
	err := database.DB.Collection("roles").FindOne(ctx, filter).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Returns the role with provided role name
func GetRoleByName(name string) (*model.Role, error) {
	filter := bson.M{"name": name}
	return GetRole(filter)
}

// Return the role with provided role tag
func GetRoleByTag(tag string) (*model.Role, error) {
	filter := bson.M{"tag": tag}
	return GetRole(filter)
}

// Insert role with provided Parameters in DB
func CreateRole(r *model.Role) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("roles").InsertOne(ctx, r)
}

// Update role with provided parameters
func UpdateRole(tag string, input *model.Role) (*mongo.UpdateResult, error) {
	filter := bson.M{"tag": tag}
	update := bson.D{
		{Key: "$set", Value: *input},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("roles").UpdateOne(ctx, filter, update)
}

// Delete role with provided filter in DB
func DeleteRole(tag string) (*mongo.DeleteResult, error) {
	// Delete role from all users whose `roles` field contain it
	users, err := GetUsers(bson.M{"roles": tag})
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err != mongo.ErrNoDocuments {
		for _, u := range users {
			DeleteRoleFromUser(tag, u)
		}
	}
	// Delete role from all content type permissions. It is possible, that one ore more permission have no roles left after.
	allContentTypes, err := GetContentTypes(bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	for _, ct := range allContentTypes {
		DeleteRoleFromPermissions(tag, ct)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"tag": tag}
	return database.DB.Collection("roles").DeleteOne(ctx, filter)
}

// Return true if the a role with given string role name exists
func IsValidRole(role model.Role) bool {
	if _, err := GetRoleByTag(role.Tag); err != nil {
		return false
	} else {
		return true
	}
}

// Returns slice of role names of provided role tags
func GetRoleNamesByTag(roleTags []string) ([]string, error) {
	var output []string
	for _, tag := range roleTags {
		rObj, err := GetRoleByTag(tag)
		if err != nil {
			return nil, err
		}
		output = append(output, rObj.Name)
	}
	return output, nil
}

func GetRoleTagsFromRoleSlice(roles []model.Role) []string {
	roleTags := make([]string, 0)
	for _, r := range roles {
		roleTags = append(roleTags, r.Tag)
	}
	return roleTags
}
