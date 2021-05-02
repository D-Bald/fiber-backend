package handler

import (
	"regexp"
	"time"

	"github.com/D-Bald/fiber-backend/config"
	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers get all Users
func GetAllUsers(c *fiber.Ctx) error {
	users, err := controller.GetUsers(bson.M{})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No User found", "user": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "user": users})
}

// GetUser get a user by ID
// Deprecated: Use GetContent with id in route parameter instead
func GetUserById(c *fiber.Ctx) error {
	user, err := controller.GetUserById(c.Params("id"))
	if err != nil || user.Username == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found", "user": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "user": user})
}

// Query users with filter provided in query params
func GetUsers(c *fiber.Ctx) error {
	re := regexp.MustCompile(`[a-z]+=[a-zA-Z0-9\%]+`)
	filterString := re.FindAllString(c.Params("*"), -1)
	filter := make(map[string]interface{})
	for _, v := range filterString {
		v = regexp.MustCompile(`%20`).ReplaceAllString(v, " ")
		s := regexp.MustCompile(`=`).Split(v, -1)
		if s[0] == `id` {
			cID, err := primitive.ObjectIDFromHex(s[1])
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "user": err.Error()})
			}
			filter["_id"] = cID
		} else {
			filter[s[0]] = s[1]
		}
	}

	result, err := controller.GetUsers(filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "user": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "user": result})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	user := new(model.User)

	// Parse input
	if err := c.BodyParser(user); err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "user": err.Error()})
	}

	// Check if already exists
	if u, _ := controller.GetUserByUsername(user.Username); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken", "user": nil})
	}
	if u, _ := controller.GetUserByEmail(user.Email); u != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists", "user": nil})
	}

	// Insert in DB
	if _, err := controller.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create user", "user": err.Error()})
	}

	// Token for response
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID.Hex()
	claims["admin"] = hasRole(user.Roles, "admin")
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// User for response
	newUser := model.UserOutput{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Names:    user.Names,
		Roles:    user.Roles,
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "token": t, "user ": newUser})
}

// Update user with parameters from request body
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !isValidToken(token, id) && !isAdminToken(token) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "result": nil})
	}

	uui := new(model.UpdateUserInput)
	if err := c.BodyParser(uui); err != nil || uui == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	if uui.Username != "" {
		if u, _ := controller.GetUserByUsername(uui.Username); u != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username already taken", "result": nil})
		}
	}
	if uui.Email != "" {
		if u, _ := controller.GetUserByEmail(uui.Email); u != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "User with given Email already exists", "result": nil})
		}
	}
	if uui.Roles != nil {
		if !isAdminToken(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Admin rights required to update user roles", "result": nil})
		}
	}

	result, err := controller.UpdateUser(id, uui)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update User", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "result": result})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password" xml:"password" form:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !isValidToken(token, id) && !isAdminToken(token) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "result": nil})
	}

	if !isValidUser(id, pi.Password) && !isAdminToken(token) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Not valid user", "result": nil})
	}

	result, err := controller.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete User", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "result": result})
}

// Validators

// Checks if the user_id claim of the token matches the id of the target user
func isValidToken(t *jwt.Token, id string) bool {
	return t.Claims.(jwt.MapClaims)["user_id"] == id
}

// Checks if the role claim of the token is `admin`
func isAdminToken(t *jwt.Token) bool {
	return t.Claims.(jwt.MapClaims)["admin"].(bool)
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
