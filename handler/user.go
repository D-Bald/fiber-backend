package handler

import (
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"

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
	newUser := model.UserOutput{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Names:    user.Names,
		Roles:    user.Roles,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// Update user with parameters from request body
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !isValidToken(token, id) && !isAdminToken(token) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	uui := new(model.UpdateUserInput)
	if err := c.BodyParser(uui); err != nil || uui == nil {
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
	if uui.Roles != nil {
		if !isAdminToken(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Admin rights required to update user roles", "data": nil})
		}
	}

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

	if !isValidToken(token, id) && !isAdminToken(token) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !isValidUser(id, pi.Password) && !isAdminToken(token) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	}

	result, err := controller.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete User", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": result})
}

// Validators

// Checks if the user_id claim of the token matches the id of the target user
func isValidToken(t *jwt.Token, id string) bool {
	return t.Claims.(jwt.MapClaims)["user_id"] == id
}

// Checks if the role claim of the token is `admin`
func isAdminToken(t *jwt.Token) bool {
	return t.Claims.(jwt.MapClaims)["admin"] == true
}

// hasRole takes a string slice of roles and looks for an element in it. If found it will
// return true, otherwise it will return false.
func hasRole(slice []string, role string) bool {
	for _, item := range slice {
		if item == role {
			return true
		}
	}
	return false
}

// Checks if the user exists in the DB and if the provided password matches the saved one
func isValidUser(id string, p string) bool {
	user, err := controller.GetUserById(id)
	if err != nil || user.Username == "" {
		return false
	}
	pw, err := controller.GetUserPasswordHash(id)
	if err != nil {
		return false
	}
	if !checkPasswordHash(p, pw) {
		return false
	}
	return true
}

// CheckPasswordHash compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
