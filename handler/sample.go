package handler

import (
	"context"

	"github.com/D-Bald/fiber-backend/database"
	"github.com/gofiber/fiber/v2"
)

type Sample struct {
	TestField1 string `bson:"testfield1" json:"testfield1" xml:"testfield1" form:"testfield1"`
	TestField2 string `bson:"testfield2" json:"testfield2" xml:"testfield2" form:"testfield2"`
}

func CreateSample(c *fiber.Ctx) error {
	sample := new(Sample)
	if err := c.BodyParser(sample); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create content", "data": err.Error()})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := database.Mg.Db.Collection("samples").InsertOne(ctx, &sample); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create sample", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": sample})
}

func GetSample(c *fiber.Ctx) error {
	sample := Sample{"Hello", "World"}
	return c.JSON(fiber.Map{"status": "success", "message": "GET Handler", "data": sample})
}
