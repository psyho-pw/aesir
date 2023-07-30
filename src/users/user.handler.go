package users

import "C"
import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
	"strconv"
)

type Handler interface {
	FindMany(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
}

type userHandler struct {
	service Service
}

func NewUserHandler(service Service) Handler {
	return &userHandler{service: service}
}

var SetHandler = wire.NewSet(NewUserHandler)

func (handler userHandler) FindMany(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	result, err := handler.service.WithTx(tx).FindMany()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler userHandler) FindOne(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id, parseErr := strconv.Atoi(c.Params("id"))
	if parseErr != nil {
		return errors.HandleParseError(c, parseErr)
	}

	result, err := handler.service.WithTx(tx).FindOne(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
