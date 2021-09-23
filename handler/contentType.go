package handler

import (
	"fmt"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Content Types
func GetAllContentTypes(c *fiber.Ctx) error {
	contentTypes, err := controller.GetContentTypes(nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "contenttype": err.Error()})
	}

	// Return a subset of fields in readable format
	result := make([]contentTypeOutput, 0)
	for _, ct := range contentTypes {
		out, err := toContentTypeOutput(ct)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on parsing permissions", "data": err.Error()})
		}
		result = append(result, *out)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Content Types", "contenttype": result})
}

// GetContentType query contenttypes by ID
func GetContentType(c *fiber.Ctx) error {
	ct, err := controller.GetContentTypeById(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "contenttype": err.Error()})

	}
	// Return a subset of fields in readable format
	ctOutput, err := toContentTypeOutput(ct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on parsing permissions", "user": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content Type found", "contenttype": ctOutput})
}

// CreateContentType
func CreateContentType(c *fiber.Ctx) error {
	ctInput := new(model.ContentTypeInput)
	// Parse input
	if err := c.BodyParser(ctInput); err != nil || ctInput.TypeName == "" || ctInput.Collection == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Review your input: 'typename' and 'collection' required", "contenttype": err.Error()})
	}

	// Check if content type already exists
	checkTypeName, _ := controller.GetContentTypeByTypeName(ctInput.TypeName)
	if checkTypeName != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "contenttype": nil})
	}
	checkCollection, _ := controller.GetContentTypeByCollection(ctInput.Collection)
	if checkCollection != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Collection name already in use for content type: %s", checkCollection.TypeName), "contenttype": nil})
	}

	// Check if all roles are valid
	if ctInput.Permissions != nil {
		for _, val := range ctInput.Permissions {
			for _, role := range val {
				if !controller.IsValidRole(role) {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role not found: %s", role), "result": nil})
				}
			}
		}
	}

	// Create actual content type object
	ct := model.ContentType{
		TypeName:    ctInput.TypeName,
		Collection:  ctInput.Collection,
		Permissions: ctInput.Permissions,
		FieldSchema: ctInput.FieldSchema,
	}

	// Insert in DB
	if _, err := controller.CreateContentType(&ct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content Type", "contenttype": err.Error()})
	}

	// Return a subset of fields in readable format
	ctOutput, err := toContentTypeOutput(&ct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on parsing permissions", "user": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created Content Type", "contenttype": ctOutput})
}

// Update content type with parameters from request body
func UpdateContentType(c *fiber.Ctx) error {
	id := c.Params("id")

	ctInput := new(model.ContentTypeInput)
	// Parse input
	if err := c.BodyParser(ctInput); err != nil || ctInput == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	// Checks if content type already exists
	if ctInput.TypeName != "" {
		checkTypeName, _ := controller.GetContentTypeByTypeName(ctInput.TypeName)
		if checkTypeName != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "contenttype": nil})
		}
	}
	if ctInput.Collection != "" {
		checkCollection, _ := controller.GetContentTypeByCollection(ctInput.Collection)
		if checkCollection != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Collection name already in use for content type: %s", checkCollection.TypeName), "contenttype": nil})
		}
	}

	// Checks, if all role are valid
	if ctInput.Permissions != nil {
		for _, val := range ctInput.Permissions {
			for _, role := range val {
				if !controller.IsValidRole(role) {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role not found: %s", role), "result": nil})
				}
			}
		}
	}
	result, err := controller.UpdateContentType(id, ctInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update Content Type", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content Type successfully updated", "result": result})
}

// DeleteContentType
func DeleteContentType(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if content type with given id exists
	if _, err := controller.GetContentTypeById(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content Type not found", "result": err.Error()})
	}

	// Delete in DB
	result, err := controller.DeleteContentType(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content Type", "result": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content Type successfully deleted", "result": result})
}

// Fields that are returned on GET methods (metadata omitted)
type contentTypeOutput struct {
	ID          string                  `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	TypeName    string                  `bson:"typename" json:"typename" xml:"typename" form:"typename"`
	Collection  string                  `bson:"collection" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]model.Role `bson:"permissions" json:"permissions" xml:"permissions" form:"permissions"`
	FieldSchema map[string]interface{}  `bson:"field_schema" json:"field_schema" xml:"field_schema" form:"field_schema"`
}

// Make ContentTypeOutput from ContentType
func toContentTypeOutput(contentType *model.ContentType) (*contentTypeOutput, error) {
	ctOutput := new(contentTypeOutput)
	ctOutput.ID = contentType.ID
	ctOutput.TypeName = contentType.TypeName
	ctOutput.Collection = contentType.Collection
	ctOutput.Permissions = contentType.Permissions
	ctOutput.FieldSchema = contentType.FieldSchema
	return ctOutput, nil
}
