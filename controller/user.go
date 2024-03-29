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

// Initialize Users with an admin user
func InitAdminUser() error {
	adminRole, err := GetRoleByTag("admin")
	if err != nil {
		return err
	}
	// Checks, if a user with a role with tag 'admin' exists, if not, create one
	_, err = GetUser(bson.M{"roles": adminRole.ID})
	if err != nil && err == mongo.ErrNoDocuments {
		hash, err := hashPassword(config.Config("FIBER_ADMIN_PASSWORD"))
		if err != nil {
			return err
		}

		// add admin and default roles to admin user
		var roles []primitive.ObjectID
		roles = append(roles, adminRole.ID)
		if userRole, err := GetRoleByTag("default"); err != nil {
			return err
		} else {
			roles = append(roles, userRole.ID)
		}

		adminUser := bson.D{
			{Key: "username", Value: "adminUser"},
			{Key: "email", Value: "admin@sample.com"},
			{Key: "password", Value: hash},
			{Key: "names", Value: "admin user"},
			{Key: "roles", Value: roles},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := database.DB.Collection("users").InsertOne(ctx, adminUser); err != nil {
			return err
		}
	}

	err = nil
	return err
}

// Return all users from DB with provided Filter
func GetUsers(filter interface{}) ([]*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var users []*model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var u model.User
		err := cursor.Decode(&u)
		if err != nil {
			return nil, err
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

// Return a single user that matches the id input
func GetUserById(id string) (*model.User, error) {
	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": uID}
	return GetUser(filter)
}

// Return a single user that matches the email input
func GetUserByEmail(e string) (*model.User, error) {
	filter := bson.M{"email": e}
	return GetUser(filter)
}

// Return a single user that matches the username input
func GetUserByUsername(u string) (*model.User, error) {
	filter := bson.M{"username": u}
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
func UpdateUser(id string, input *model.UserUpdate) (*mongo.UpdateResult, error) {
	// Struct similar to `UserUpdate` but with ObjectIDs of roles instead of string role names
	type mongoUserUpdate struct {
		Username string               `bson:"username,omitempty"`
		Email    string               `bson:"email,omitempty"`
		Password string               `bson:"password,omitempty"`
		Names    string               `bson:"names,omitempty"`
		Roles    []primitive.ObjectID `bson:"roles,omitempty"`
	}

	// create Object to add ObjectIDs as Roles
	userUpdate := mongoUserUpdate{
		Username: input.Username,
		Email:    input.Email,
		Password: "",
		Names:    input.Names,
		Roles:    make([]primitive.ObjectID, 0),
	}

	// Hash the password before updating the user
	if input.Password != "" {
		hash, err := hashPassword(input.Password)
		if err != nil {
			return new(mongo.UpdateResult), err
		}
		userUpdate.Password = hash
	}

	// Parse role name strings to role ObjectIDs
	for _, r := range input.Roles {
		rObj, err := GetRoleByName(r)
		if err != nil {
			return new(mongo.UpdateResult), err
		}
		userUpdate.Roles = append(userUpdate.Roles, rObj.ID)
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return new(mongo.UpdateResult), err
	}
	// Update user with provided ID and sets field value for `updatet_at`
	filter := bson.M{"_id": userID}
	update := bson.D{
		{Key: "$set", Value: userUpdate},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
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
	filter := bson.M{"_id": uID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("users").DeleteOne(ctx, filter)
}

// Delete only one role from user.
func DeleteRoleFromUser(rID primitive.ObjectID, user *model.User) (*mongo.UpdateResult, error) {
	roles := make([]primitive.ObjectID, 0)
	for _, r := range user.Roles {
		if r != rID {
			roles = append(roles, r)
		}
	}
	// Update user with provided ID and sets field value for `updatet_at`
	filter := bson.M{"_id": user.ID}
	update := bson.D{
		{Key: "$set", Value: bson.M{"roles": roles}},
		{Key: "$currentDate", Value: bson.M{
			"updated_at": true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return database.DB.Collection("users").UpdateOne(ctx, filter, update)
}

// Hashes password string with bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
