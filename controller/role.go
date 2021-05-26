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

// Init roles: 'admin' and 'user'
var (
	admin = bson.D{
		{Key: "role", Value: "admin"},
		{Key: "weight", Value: 1000},
	}
	user = bson.D{
		{Key: "role", Value: "user"},
		{Key: "weight", Value: 0},
	}
)

// Initialize collection roles with 'admin' and 'user'
func InitRoles() error {
	_, err := GetRole("user")
	if err != nil && err == mongo.ErrNoDocuments {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := database.DB.Collection("roles").InsertOne(ctx, user); err != nil {
			return err
		}
	}
	_, err = GetRole("admin")
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

// Returns the role with provided role name
func GetRole(role string) (*model.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"role": role}
	var r *model.Role
	err := database.DB.Collection("roles").FindOne(ctx, filter).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Returns the role Object with provided ID
func GetRoleById(id string) (*model.Role, error) {
	rID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": rID}
	var r *model.Role
	err = database.DB.Collection("roles").FindOne(ctx, filter).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Insert role with provided Parameters in DB
func CreateRole(r *model.Role) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("roles").InsertOne(ctx, r)
}

// Update role with provided parameters
func UpdateRole(input *model.Role) (*mongo.UpdateResult, error) {
	filter := bson.M{"role": input.Role}
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
func DeleteRole(filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("roles").DeleteOne(ctx, filter)
}

// Return true if the a role with given string role name exists
func IsValidRole(role string) bool {
	if _, err := GetRole(role); err != nil {
		return false
	} else {
		return true
	}
}

// Returns slice of role names of provided role ObjectsIDs
func GetRoleNames(roleIDs []primitive.ObjectID) ([]string, error) {
	var output []string
	for _, r := range roleIDs {
		rObj, err := GetRoleById(r.Hex())
		if err != nil {
			return nil, err
		}
		output = append(output, rObj.Role)
	}
	return output, nil
}
