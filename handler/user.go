package handler

import (
	"fmt"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers get all Users
func GetUsers(c *fiber.Ctx) error {
	users, err := controller.GetUsers(bson.D{})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No User found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	user, err := controller.GetUserById(c.Params("id"))
	if err != nil || user.Username == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		ID       primitive.ObjectID `json:"id"`
		Username string             `json:"username"`
		Email    string             `json:"email"`
		Role     string             `json:"role"`
	}

	user := new(model.User)

	// Parse input
	if err := c.BodyParser(user); err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	// Check if already exists
	if u, _ := controller.GetUserByUsername(user.Username); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken", "data": nil})
	}
	if u, _ := controller.GetUserByEmail(user.Email); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists", "data": nil})
	}

	// Insert in DB
	if _, err := controller.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err.Error()})
	}

	// Response
	newUser := NewUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// Update user with parameters from request body
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !isValidToken(token, id) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	uui := new(model.UpdateUserInput)
	if err := c.BodyParser(uui); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	if uui.Username != "" {
		if u, _ := controller.GetUserByUsername(uui.Username); u != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken", "data": nil})
		}
	}
	if uui.Email != "" {
		if u, _ := controller.GetUserByEmail(uui.Email); u != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists", "data": nil})
		}
	}
	// if uui.Role != "" {
	// 	// CHECK FOR ADMIN CLAIM HERE
	// }

	result, err := controller.UpdateUser(id, uui)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update User", "data": err.Error()})
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

	if !isValidToken(token, id) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !isValidUser(id, pi.Password) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	}

	result, err := controller.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete User", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": result})
}

// Validators

func isValidToken(t *jwt.Token, id string) bool {

	claims := t.Claims.(jwt.MapClaims)
	uid := claims["user_id"]

	if uid != id {
		fmt.Println("ID anders als im Claim")
		return false
	}

	return true
}

func isValidUser(id string, p string) bool {
	user, err := controller.GetUserById(id)
	if err != nil || user.Username == "" {
		return false
	}
	if !checkPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// CheckPasswordHash compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
