package controller

import (
	"context"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Return all users from DB with given Filter
func GetUsers(filter interface{}) ([]*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var users []*model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("users").Find(ctx, filter)
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var u model.User
		err := cursor.Decode(&u)
		if err != nil {
			return users, err
		}

		users = append(users, &u)
	}

	if err := cursor.Err(); err != nil {
		return users, err
	}

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

// Return a single user that matches the filter
func GetUser(filter interface{}) (*model.User, error) {
	var user *model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(id string) (*model.User, error) {
	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: uID}}
	return GetUser(filter)
}

func GetUserByEmail(e string) (*model.User, error) {
	filter := bson.D{{Key: "email", Value: e}}
	return GetUser(filter)
}

func GetUserByUsername(u string) (*model.User, error) {
	filter := bson.D{{Key: "username", Value: u}}
	return GetUser(filter)
}

// Insert user with given Parameters in DB
func CreateUser(user *model.User) (*mongo.InsertOneResult, error) {
	// Initialize metadata
	user.Init()

	// Hash the password before saving the user
	hash, err := hashPassword(user.Password)
	if err != nil {
		return new(mongo.InsertOneResult), err
	}
	user.Password = hash

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return database.DB.Collection("users").InsertOne(ctx, user)
}

// Update user with given Parameters in DB
func UpdateUser(id string, input *model.UpdateUserInput) (*mongo.UpdateResult, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return new(mongo.UpdateResult), err
	}

	// Hash the password before updating the user
	if input.Password != "" {
		hash, err := hashPassword(input.Password)
		if err != nil {
			return new(mongo.UpdateResult), err
		}
		input.Password = hash
	}

	// Update user with given ID: sets field values for "names" and "updatet_at"
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$set", Value: *input},
		{Key: "$currentDate", Value: bson.D{
			{Key: "updated_at", Value: true},
		}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("users").UpdateOne(ctx, filter, update)
}

// Delete user with given ID in DB
func DeleteUser(id string) (*mongo.DeleteResult, error) {
	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: uID}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("users").DeleteOne(ctx, filter)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
