package handler

import (
	"regexp"

	"github.com/D-Bald/fiber-backend/controller"
	"github.com/D-Bald/fiber-backend/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Return all content entries in given collection
func GetAllContentEntries(c *fiber.Ctx) error {
	coll := c.Params("content")

	result, err := controller.GetContent(coll, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Content Entries", "data": result})
}

// Query content by Param 'id'
// Deprecated: Use GetContent with id in route parameter instead
func GetContentById(c *fiber.Ctx) error {
	coll := c.Params("content")
	content, err := controller.GetContentById(coll, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": content})
}

// Query content entries with filter provided in query params
// Solution with regular expressions
// Every Value is Parsed as string type.
func GetContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	re := regexp.MustCompile(`[a-z\_]+=[a-zA-Z0-9\%]+`)
	filterString := re.FindAllString(c.Params("*"), -1)
	filter := make(map[string]interface{})
	for _, v := range filterString {
		v = regexp.MustCompile(`%20`).ReplaceAllString(v, " ")
		s := regexp.MustCompile(`=`).Split(v, -1)

		// check, if the current param is boolean to parse it in a boolean output if this is the case
		boolMatch, boolOutput, err := parseBoolean(s[1])
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
		}
		switch {
		case s[0] == `id`:
			cID, err := primitive.ObjectIDFromHex(s[1])
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
			}
			filter[s[0]] = cID

		case boolMatch:
			filter[s[0]] = boolOutput
		default:
			filter[s[0]] = s[1]
		}
	}

	result, err := controller.GetContent(coll, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": result})
}

// // Query content entries with filter provided in query params
// // Solution with QueryParser
// // DOES NOT WORK BECAUSE FIELDS ARE INITIALIZED WITH NIL VALUES WHICH DOES NOT MATCH ANY DOCUMENT => USE c.Query("KEY") for single keys?
// func GetContent(c *fiber.Ctx) error {
// 	coll := c.Params("content")
// 	filter := new(model.Content)

// 	// parse input
// 	if err := c.QueryParser(filter); err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
// 	}

// 	// get content from DB
// 	result, err := controller.GetContent(coll, filter)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No match found", "data": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Content found", "data": result})
// }

// CreateContent new content
// Collection is created by mongoDB automatically on first insert call
func CreateContent(c *fiber.Ctx) error {
	content := new(model.Content)
	// Parse input
	if err := c.BodyParser(content); err != nil || content.Title == "" || content.Fields == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create content", "data": err.Error()})
	}

	// Get collection from route params
	coll := c.Params("content")

	if _, err := controller.CreateContent(coll, content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create Content", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": content})
}

// Update content entry with parameters from request body
// Update user with parameters from request body
func UpdateContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")

	uci := new(model.UpdateContentInput)
	if err := c.BodyParser(uci); err != nil || uci == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	result, err := controller.UpdateContent(coll, id, uci)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update content entry", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully updated", "data": result})
}

// DeleteContent delete content
func DeleteContent(c *fiber.Ctx) error {
	coll := c.Params("content")
	id := c.Params("id")
	if _, err := controller.GetContentById(coll, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Content not found", "data": nil})
	}

	result, err := controller.DeleteContent(coll, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not delete Content", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Content successfully deleted", "data": result})
}

// if input s is the string of a boolean keyword, return that keyword
// false is default, but if isMatch is false, the output is not valid.
func parseBoolean(s string) (isMatch bool, output bool, err error) {
	t, err := regexp.MatchString(`((t|T)rue)`, s)
	if err != nil {
		return false, false, err
	}
	if t {
		return true, true, nil
	}
	f, err := regexp.MatchString(`(f|F)alse`, s)
	if err != nil {
		return false, false, err
	}
	if f {
		return true, false, nil
	}
	return false, false, err
}
