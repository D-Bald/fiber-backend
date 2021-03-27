package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "github.com/dgrijalva/jwt-go" <- Nicht kompatibel mit "github.com/gofiber/jwt/v2", was hier in Tokens aus dem fiber Context verwendet wird.
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	if uid != n {
		fmt.Println("ID anders als im Claim")
		return false
	}

	return true
}

func validUser(id string, p string) bool {
	userID, _ := primitive.ObjectIDFromHex(id)
	user, err := findUser(bson.M{"_id": userID})
	if err != nil || user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// Filter Users with given Filter
func filterUsers(filter interface{}) ([]*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var users []*model.User
	ctx := context.TODO()
	cursor, err := database.Mg.Db.Collection("User").Find(ctx, filter)
	if err != nil {
		return users, err
	}

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

	// once exhausted, close the cursor
	cursor.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

// Find a single User by filter
func findUser(filter interface{}) (*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var user *model.User
	ctx := context.TODO()
	err := database.Mg.Db.Collection("User").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers get all Users in DB
func GetUsers(c *fiber.Ctx) error {
	users, err := filterUsers(bson.M{})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No User found.", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on User ID", "data": nil})
	}
	user, err := findUser(bson.M{"_id": userID})
	if err != nil || user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		ID       primitive.ObjectID `json:"id" xml:"id" form:"id"`
		Username string             `json:"username" xml:"username" form:"username"`
		Email    string             `json:"email" xml:"email" form:"email"`
	}

	user := new(model.User)

	if err := c.BodyParser(user); err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	if _, err := findUser(bson.M{"username": user.Username}); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken"})
	}
	if _, err := findUser(bson.M{"email": user.Email}); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists"})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}

	user.Password = hash
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	col := database.Mg.Db.Collection("users")
	if _, err := col.InsertOne(context.TODO(), &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names" xml:"names" form:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on User ID", "data": nil})
	}

	// Update User with given ID: sets Field Values for "names" and "updatet_at"
	filter := bson.M{"_id": userID}
	update := bson.D{
		{"$set", bson.D{
			{"names", uui.Names},
		}},
		{"$currentDate", bson.D{
			{"updated_at", true},
		}},
	}
	result, err := database.Mg.Db.Collection("users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update User", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": result})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password" xml:"password" form:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !validUser(id, pi.Password) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	}
	userID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userID}
	result, err := database.Mg.Db.Collection("users").DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete User", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": result})
}
