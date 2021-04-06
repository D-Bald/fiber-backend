package controller

import (
	"context"
	"time"

	"github.com/D-Bald/fiber-backend/config"
	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Initialize Collection Users with a admin user
func hashedAdminPassword() (string, error) {
	hash, err := hashPassword(config.Config("ADMIN_PASSWORD"))
	if err != nil {
		return bson.TypeNull.String(), err
	}
	return hash, nil
}

func InitAdminUser() error {
	_, err := GetUsers(bson.D{{Key: "roles", Value: "admin"}})
	if err != nil {
		hash, err := hashedAdminPassword()
		if err != nil {
			return err
		}
		adminUser := bson.D{
			{Key: "username", Value: "adminUser"},
			{Key: "email", Value: "admin@sample.com"},
			{Key: "password", Value: hash},
			{Key: "names", Value: "admin user"},
			{Key: "roles", Value: bson.A{"user", "admin"}},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := database.DB.Collection("users").InsertOne(ctx, adminUser); err != nil {
			return err
		}
	}
	return err
}

// Return all users from DB with provided Filter
func GetUsers(filter interface{}) ([]*model.UserOutput, error) {
	// A slice of tasks for storing the decoded documents
	var users []*model.UserOutput

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("users").Find(ctx, filter)
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var u model.UserOutput
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
func GetUser(filter interface{}) (*model.UserOutput, error) {
	var user *model.UserOutput

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Return a single user that matches the id input
func GetUserById(id string) (*model.UserOutput, error) {
	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: uID}}
	return GetUser(filter)
}

// Return a single user that matches the email input
func GetUserByEmail(e string) (*model.UserOutput, error) {
	filter := bson.D{{Key: "email", Value: e}}
	return GetUser(filter)
}

// Return a single user that matches the username input
func GetUserByUsername(u string) (*model.UserOutput, error) {
	filter := bson.D{{Key: "username", Value: u}}
	return GetUser(filter)
}

// Return hashed Password of User with provided ID
func GetUserPasswordHash(id string) (string, error) {
	type userPassword struct {
		Password string `bson:"password"`
	}
	var userPW *userPassword

	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return bson.TypeNull.String(), err
	}
	filter := bson.M{"_id": uID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.DB.Collection("users").FindOne(ctx, filter).Decode(&userPW)
	if err != nil {
		return bson.TypeNull.String(), err
	}

	return userPW.Password, nil
}

// Insert user with provided Parameters in DB
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

// Update user with provided Parameters in DB
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

	// Update user with provided ID: sets field values for "names" and "updatet_at"
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

// Delete user with provided ID in DB
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
