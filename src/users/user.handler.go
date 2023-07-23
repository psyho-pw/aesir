package users

import "C"
import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"strconv"
)

type UserHandler interface {
	FindMany(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
}

type userHandler struct {
	service UserService
}

func NewUserHandler(service UserService) UserHandler {
	return &userHandler{service: service}
}

var SetHandler = wire.NewSet(NewUserHandler)

func (handler userHandler) FindMany(c *fiber.Ctx) error {
	result, err := handler.service.FindMany()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler userHandler) FindOne(c *fiber.Ctx) error {
	id, parseErr := strconv.Atoi(c.Params("id"))
	if parseErr != nil {
		return errors.HandleParseError(c, parseErr)
	}

	result, err := handler.service.FindOne(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
