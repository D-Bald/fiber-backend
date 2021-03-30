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

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// private funcs
// Auth and Validation
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
	user, err := getUserById(id)
	if err != nil || user.Username == "" {
		return false
	}
	if !checkPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// Getters for easy DB lookups
// Return Users from DB with given Filter
func getUsers(filter interface{}) ([]*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var users []*model.User
	ctx := context.TODO()
	cursor, err := database.Mg.Db.Collection("users").Find(ctx, filter)
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

// Return a single User from DB by filter
func getUser(filter interface{}) (*model.User, error) {
	// A slice of tasks for storing the decoded documents
	var user *model.User
	ctx := context.TODO()
	err := database.Mg.Db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUserById(id string) (*model.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: userID}}
	user, err := getUser(filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByEmail(e string) (*model.User, error) {
	filter := bson.D{{Key: "email", Value: e}}
	user, err := getUser(filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	filter := bson.D{{Key: "username", Value: u}}
	user, err := getUser(filter)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// public funcs
// actual Handlers
// GetUsers get all Users in DB
func GetUsers(c *fiber.Ctx) error {
	users, err := getUsers(bson.D{})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No User found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	user, err := getUserById(c.Params("id"))
	if err != nil || user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		ID       primitive.ObjectID `json:"id"`
		Username string             `json:"username"`
		Email    string             `json:"email"`
	}

	user := new(model.User)

	if err := c.BodyParser(user); err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	if u, _ := getUserByUsername(user.Username); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken", "data": nil})
	}
	if u, _ := getUserByEmail(user.Email); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists", "data": nil})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not hash password", "data": err.Error()})

	}

	user.Password = hash
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if _, err := database.Mg.Db.Collection("users").InsertOne(context.TODO(), &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err.Error()})
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
		Names string `json:"names" xml:"names" form:"names"` // erweitern, sondass auch andere Felder geupdatet werden können
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error on User ID", "data": err.Error()})
	}

	// Update User with given ID: sets Field Values for "names" and "updatet_at"
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "names", Value: uui.Names},
		}},
		{Key: "$currentDate", Value: bson.D{
			{Key: "updated_at", Value: true},
		}},
	}
	result, err := database.Mg.Db.Collection("users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update User", "data": err.Error()})
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
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
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
	filter := bson.D{{Key: "_id", Value: userID}}
	result, err := database.Mg.Db.Collection("users").DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete User", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": result})
}
