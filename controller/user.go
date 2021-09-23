package controller

import (
	"context"
	"reflect"
	"time"

	"github.com/D-Bald/fiber-backend/config"
	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Content struct for Mongo Driver
// Shadowing the ID and Permissions fields of model.User with ObjectID for ID and just the string slice of Role tags. Database Entry need not to be updated on role-update then.
type UserMongo struct {
	User  *model.User
	ID    primitive.ObjectID `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	Roles []string           `bson:"roles" json:"roles" xml:"roles" form:"roles" query:"roles"`
}

// Return the User Struct and change ID to string and Permissions to roles
func (u *UserMongo) toUser() (*model.User, error) {
	user := u.User
	user.ID = u.ID.Hex()
	for _, tag := range u.Roles {
		role, err := GetRoleByTag(tag)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, *role)
	}
	return user, nil
}

// Initialize Users with an admin user
func InitAdminUser() error {
	// Checks, if a user with a role with tag 'admin' exists, if not, create one
	_, err := GetUser(bson.M{"roles": "admin"})
	if err != nil && err == mongo.ErrNoDocuments {
		hash, err := hashPassword(config.Config("FIBER_ADMIN_PASSWORD"))
		if err != nil {
			return err
		}

		// add admin and default roles to admin user
		var roles []string
		roles = append(roles, "admin")
		roles = append(roles, "default")

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
	var result []*model.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("users").Find(ctx, filter)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var uMongo UserMongo
		err := cursor.Decode(&uMongo)
		if err != nil {
			return result, err
		}
		user, err := uMongo.toUser()
		if err != nil {
			return result, err
		}
		result = append(result, user)
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}

// Return a single user that matches the filter
func GetUser(input interface{}) (*model.User, error) {
	// set `nil` for empty values
	v := reflect.ValueOf(input)
	filter := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			// Differentiates between different fields of the struct specified by their json flag
			switch v.Type().Field(i).Tag.Get("json") {
			// Parses ID manually to ObjectID and add it to filter
			case "_id":
				uID, err := primitive.ObjectIDFromHex(v.Field(i).String())
				if err != nil {
					return nil, err
				}
				filter["_id"] = uID

			// Parses roles manually to Tags and add it to filter
			case "roles":
				roleTags := GetRoleTagsFromRoleSlice(v.Field(i).Interface().([]model.Role))

				// Checks if the query slice contains only one value. If so, add this value; Add a slice otherwise
				if len(roleTags) == 1 {
					filter["roles"] = roleTags[0]
				} else {
					filter["roles"] = roleTags
				}
			// add any other parameter to the filter
			default:
				filter[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Interface()
			}
		}

		// Check for boolean types, because the zero value of this type `false` can be relevant for queries
		if v.Type().Field(i).Type.Kind() == reflect.Bool {
			filter[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Interface()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userMongo *UserMongo

	err := database.DB.Collection("users").FindOne(ctx, filter).Decode(&userMongo)
	if err != nil {
		return nil, err
	}

	result, err := userMongo.toUser()
	if err != nil {
		return nil, err
	}

	return result, nil
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
func UpdateUser(id string, input *model.UserInput) (*mongo.UpdateResult, error) {
	userMongo := new(UserMongo)

	userMongo.User.Username = input.Username
	userMongo.User.Email = input.Email
	userMongo.User.Names = input.Names
	// Hash the password before updating the user
	if input.Password != "" {
		hash, err := hashPassword(input.Password)
		if err != nil {
			return new(mongo.UpdateResult), err
		}
		userMongo.User.Password = hash
	}

	// only store role tags
	userMongo.Roles = GetRoleTagsFromRoleSlice(input.Roles)

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return new(mongo.UpdateResult), err
	}
	// Update user with provided ID and sets field value for `updatet_at`
	filter := bson.M{"_id": userID}
	update := bson.D{
		{Key: "$set", Value: userMongo},
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
func DeleteRoleFromUser(delRoleTag string, user *model.User) (*mongo.UpdateResult, error) {
	roles := make([]string, 0)
	for _, r := range user.Roles {
		if r.Tag != delRoleTag {
			roles = append(roles, r.Tag)
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
