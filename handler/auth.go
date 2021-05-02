package handler

import (
	"time"

	"github.com/D-Bald/fiber-backend/config"
	"github.com/D-Bald/fiber-backend/controller"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity" xml:"identity" form:"identity"`
		Password string `json:"password" xml:"password" form:"password"`
	}
	type UserData struct {
		ID       primitive.ObjectID `json:"id"`
		Username string             `json:"username"`
		Email    string             `json:"email"`
		Names    string             `json:"names"`
		Roles    []string           `json:"roles"`
	}
	var input LoginInput
	var ud UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "token": nil, "user": nil})
	}

	identity := input.Identity
	if identity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "No Identity provided on Login", "token": nil, "user": nil})
	}
	pass := input.Password

	email, _ := controller.GetUserByEmail(identity)

	user, _ := controller.GetUserByUsername(identity)

	if email == nil && user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found", "token": nil, "user": nil})
	}

	if email == nil {
		ud = UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Names:    user.Names,
			Roles:    user.Roles,
		}
	} else {
		ud = UserData{
			ID:       email.ID,
			Username: email.Username,
			Email:    email.Email,
			Names:    email.Names,
			Roles:    email.Roles,
		}
	}

	pw, err := controller.GetUserPasswordHash(ud.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not validate user", "token": nil, "user": nil})
	}

	if !checkPasswordHash(pass, pw) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "invalid password", "token": nil, "user": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID.Hex()
	claims["admin"] = hasRole(ud.Roles, "admin")
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": fiber.Map{"token": t, "user": ud}})
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "user": ud})
}
