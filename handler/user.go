package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	db := database.Mg.Db
	col := db.Collection("Users")
	var user model.User
	err = col.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&user)
	if err != nil || user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// GetUsers get all Users in DB
func GetUsers(c *fiber.Ctx) error {
	db := database.Mg.Db
	var users []model.User
	db.Find(&users)
	if users == nil || len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No User found.", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	docID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Error on User ID", "data": nil})
	}
	db := database.Mg.Db
	col := db.Collection("users")
	var user model.User
	err = col.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&user)
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

	if err := c.BodyParser(user); err != nil || user.Username == "" || user.Email == "" {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})

	}

	if u, err := getUserByUsername(user.Username); err == nil && u != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Username already taken"})
	}
	if e, err := getUserByEmail(user.Email); err == nil && e != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists"})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}

	user.Password = hash

	db := database.Mg.Db
	col := db.Collection("users")
	if _, err := col.InsertOne(context.TODO(), &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
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
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.Mg.Db
	var user model.User

	db.First(&user, id)
	user.Names = uui.Names
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password" xml:"password" form:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})

	}

	if !validUser(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})

	}

	db := database.Mg.Db
	var user model.User

	db.First(&user, id)

	db.Delete(&user)
	/*
	 * TO DO: Delete(&user) Setzt nur das Feld 'DeletedAt' auf aktuelle Zeit, behält die User aber in der Datenbank.
	 * "Gelöschte User" werden jedoch nicht mehr bei GET Anfragen ausgegeben! => Hier: Lieber ganz aus der D löschen!
	 * Zum Beispiel: nicht gorm.Model im model embedden => 'DeletedAt' nicht enthalten.
	 */
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
