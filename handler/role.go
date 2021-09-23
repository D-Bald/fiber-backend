package handler

import (
	"fmt"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Roles
func GetRoles(c *fiber.Ctx) error {
	result, err := controller.GetRoles(nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "role": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Roles", "role": result})
}

// CreateRole
func CreateRole(c *fiber.Ctx) error {
	role := new(model.Role)

	// Parse input
	if err := c.BodyParser(role); err != nil || role.Tag == "" || role.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input: 'tag' and 'name' required", "role": err.Error()})
	}

	// Check if already exists
	checkRoleTag, _ := controller.GetRoleByTag(role.Tag)
	if checkRoleTag != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role tag already in use with role name: %s", checkRoleTag.Name), "role": nil})
	}
	checkRoleName, _ := controller.GetRoleByName(role.Name)
	if checkRoleName != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role name already in use with role tag: %s", checkRoleName.Tag), "role": nil})
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
	tag := c.Params("tag")
	r := new(model.Role)
	if err := c.BodyParser(r); err != nil || r == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	// Check if role exists
	// Check if already exists
	checkRoleTag, _ := controller.GetRoleByTag(r.Tag)
	if checkRoleTag != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role tag already in use with role name: %s", checkRoleTag.Name), "result": nil})
	}
	checkRoleName, _ := controller.GetRoleByName(r.Name)
	if checkRoleName != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role name already in use with role tag: %s", checkRoleName.Tag), "result": nil})
	}

	result, err := controller.UpdateRole(tag, r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update role", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Role successfully updated", "result": result})
}

// DeleteRole delete role with provided role name
func DeleteRole(c *fiber.Ctx) error {
	tag := c.Params("tag")

	// Check if role exists
	if _, err := controller.GetRoleByTag(tag); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Role not found", "result": err.Error()})
	}

	// Delete in DB
	result, err := controller.DeleteRole(tag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete role", "result": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Role successfully deleted", "result": result})
}
