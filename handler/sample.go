package handler

import "github.com/gofiber/fiber/v2"

type Sample struct {
	TestField1 string `bson:"testfield1" json:"testfield1" xml:"testfield1" form:"testfield1"`
	TestField2 string `bson:"testfield2" json:"testfield2" xml:"testfield2" form:"testfield2"`
}

func CreateSample(c *fiber.Ctx) error {
	sample := new(Sample)
	if err := c.BodyParser(sample); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create content", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created content", "data": sample})
}

func GetSample(c *fiber.Ctx) error {
	sample := Sample{"Hello", "World"}
	return c.JSON(fiber.Map{"status": "success", "message": "GET Handler", "data": sample})
}
