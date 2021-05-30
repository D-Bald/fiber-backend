package handler

import (
	"fmt"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

// GetAll query all Content Types
func GetAllContentTypes(c *fiber.Ctx) error {
	contentTypes, err := controller.GetContentTypes(bson.M{})
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
	type newContentType struct {
		TypeName    string                 `bson:"typename" json:"typename"`
		Collection  string                 `bson:"collection" json:"collection"`
		Permissions map[string][]string    `bson:"permissions" json:"permissions"`
		FieldSchema map[string]interface{} `bson:"field_schema" json:"field_schema"`
	}

	ctInput := new(newContentType)
	// Parse input
	if err := c.BodyParser(ctInput); err != nil || ctInput.TypeName == "" || ctInput.Collection == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Review your input: 'typename' and 'collection' required", "contenttype": err.Error()})
	}

	// Check if content type already exists
	checkTypeName, _ := controller.GetContentType(bson.M{"typename": ctInput.TypeName})
	checkCollection, _ := controller.GetContentType(bson.M{"collection": ctInput.Collection})
	if checkTypeName != nil || checkCollection != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "contenttype": nil})
	}

	// Check if all roles are valid and parse them to Object IDs
	permissions := make(map[string][]primitive.ObjectID)
	if ctInput.Permissions != nil {
		for key, val := range ctInput.Permissions {
			var roleObjectIDs []primitive.ObjectID
			for _, role := range val {
				if !controller.IsValidRole(role) {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role not found: %s", role), "result": nil})
				} else {
					rObj, err := controller.GetRoleByName(role)
					if err != nil {
						return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Could not create content type", "result": nil})
					}
					roleObjectIDs = append(roleObjectIDs, rObj.ID)
				}
			}
			permissions[key] = roleObjectIDs
		}
	}

	// Create actual content type object
	ct := model.ContentType{
		TypeName:    ctInput.TypeName,
		Collection:  ctInput.Collection,
		Permissions: permissions,
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

	ctui := new(model.ContentTypeUpdate)
	// Parse input
	if err := c.BodyParser(ctui); err != nil || ctui == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	// Checks if content type already exists
	if ctui.TypeName != "" {
		checkTypeName, _ := controller.GetContentType(bson.M{"typename": ctui.TypeName})
		if checkTypeName != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "contenttype": nil})
		}
	}
	if ctui.Collection != "" {
		checkCollection, _ := controller.GetContentType(bson.M{"collection": ctui.Collection})
		if checkCollection != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Content Type already exists", "contenttype": nil})
		}
	}

	// Checks, if all role are valid
	if ctui.Permissions != nil {
		for _, val := range ctui.Permissions {
			for _, role := range val {
				if !controller.IsValidRole(role) {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": fmt.Sprintf("Role not found: %s", role), "result": nil})
				}
			}
		}
	}
	result, err := controller.UpdateContentType(id, ctui)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update User", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "result": result})
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

// Fields that are returned on GET methods (password and metadata omitted)
type contentTypeOutput struct {
	ID          primitive.ObjectID     `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	TypeName    string                 `bson:"typename" json:"typename" xml:"typename" form:"typename"`
	Collection  string                 `bson:"collection" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]string    `bson:"permissions" json:"permissions" xml:"permissions" form:"permissions"`
	FieldSchema map[string]interface{} `bson:"field_schema" json:"field_schema" xml:"field_schema" form:"field_schema"`
}

// Make ContentTypeOutput from ContentType
func toContentTypeOutput(contentType *model.ContentType) (*contentTypeOutput, error) {
	ct := new(contentTypeOutput)
	ct.ID = contentType.ID
	ct.TypeName = contentType.TypeName
	ct.Collection = contentType.Collection
	// Parse role ObjectIDs to role name strings
	permissions := make(map[string][]string)
	for key, val := range contentType.Permissions {
		roles, err := controller.GetRoleNames(val)
		if err != nil {
			return nil, err
		}
		permissions[key] = roles
	}
	ct.Permissions = permissions
	ct.FieldSchema = contentType.FieldSchema
	return ct, nil
}
