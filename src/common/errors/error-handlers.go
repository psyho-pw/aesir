package errors

import "github.com/gofiber/fiber/v2"

func HandleParseError(c *fiber.Ctx, err error) error {
	_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "cannot parse params",
	})
	return err
}
