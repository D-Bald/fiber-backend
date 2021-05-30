package handler

import (
	"fmt"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Roles
func GetRoles(c *fiber.Ctx) error {
	result, err := controller.GetRoles(bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "role": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Roles", "role": result})
}

// CreateRole
func CreateRole(c *fiber.Ctx) error {
	role := new(model.Role)

	// Parse input
	if err := c.BodyParser(role); err != nil || role.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input: 'role' required", "role": err.Error()})
	}

	if role.Weight < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input: 'weight' > 0 required", "role": nil})
	}

	// Check if already exists
	checkRole, _ := controller.GetRoleByName(role.Role)
	if checkRole != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Role already exists", "role": nil})
	}
	// Check, if weight is already taken
	checkWeight, _ := controller.GetRoles(bson.M{})
	for _, r := range checkWeight {
		if r.Weight == role.Weight {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Weight %v already taken by role: %s", role.Weight, r.Role), "role": nil})
		}
	}

	// Insert in DB
	if _, err := controller.CreateRole(role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create role", "role": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created Role", "role": role})
}

// Update role with parameters from request body
// lookup by: role
// field to update: weight
func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	r := new(model.Role)
	if err := c.BodyParser(r); err != nil || r == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	// Check if role exists
	_, err := controller.GetRoleByName(r.Role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Role not found", "result": nil})
	}
	// Check, if weight is already taken
	checkWeight, _ := controller.GetRoles(bson.M{})
	for role := 0; role < len(checkWeight); role++ {
		if checkWeight[role].Weight == r.Weight {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Weight %v already taken by role: %s", r.Weight, checkWeight[role].Role), "result": nil})
		}
	}

	result, err := controller.UpdateRole(id, r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update role", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Role successfully updated", "result": result})
}

// DeleteRole delete role with provided role name
func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if role exists
	if _, err := controller.GetRoleById(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Role not found", "result": err.Error()})
	}

	// Delete in DB
	result, err := controller.DeleteRole(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete role", "result": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Role successfully deleted", "result": result})
}
