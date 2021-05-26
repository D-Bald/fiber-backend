package handler

import (
	"reflect"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Query content entries with filter provided in query params
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	parseObject := new(model.Content)

	// parse input
	if err := c.QueryParser(parseObject); err != nil && err.Error() != "schema: converter not found for primitive.ObjectID" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
	}

	// parse ID manually because fiber's QueryParser has no converter for this type.
	if id := string(c.Request().URI().QueryArgs().Peek("id")); id != "" {
		cID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid ID", "data": err.Error()})
		}
		parseObject.ID = cID
	}

	// set `nil`for empty values
	v := reflect.ValueOf(*parseObject)
	filter := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			// If the query Field is a Slice and contains just one value, just add the single value
			if v.Field(i).Type() == reflect.SliceOf(reflect.TypeOf("")) {
				if v.Field(i).Len() == 1 {
					filter[string(v.Type().Field(i).Tag.Get("bson"))] = v.Field(i).Index(0).Interface()
				} else {
					filter[string(v.Type().Field(i).Tag.Get("bson"))] = v.Field(i).Interface()
				}
			} else {
				filter[string(v.Type().Field(i).Tag.Get("bson"))] = v.Field(i).Interface()
			}
		}
		// Check for boolean types, because the zero value of this type `false` can be relevant for queries
		if v.Type().Field(i).Type.Kind() == reflect.Bool {
			filter[string(v.Type().Field(i).Tag.Get("bson"))] = v.Field(i).Interface()
		}
	}

	// get content from DB
	result, err := controller.GetContent(coll, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "content": result})
}

// CreateContent new content
// Collection is created by mongoDB automatically on first insert call
func CreateContent(c *fiber.Ctx) error {
	content := new(model.Content)
	// Parse input
	if err := c.BodyParser(content); err != nil || content.Title == "" || content.Fields == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create content", "content": err.Error()})
	}

	// Get collection from route params
	coll := c.Params("content")

	if _, err := controller.CreateContent(coll, content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "content": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "content": content})
}

// Update content entry with parameters from request body
// Update user with parameters from request body
func UpdateContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")

	uci := new(model.ContentUpdateInput)
	if err := c.BodyParser(uci); err != nil || uci == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "result": err.Error()})
	}

	result, err := controller.UpdateContent(coll, id, uci)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update content entry", "result": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully updated", "result": result})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")
	if _, err := controller.GetContentById(coll, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "result": err.Error()})
	}

	result, err := controller.DeleteContent(coll, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "result": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "result": result})
}
