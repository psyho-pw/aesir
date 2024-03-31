package clients

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
	"strconv"
)

type Handler interface {
	CreateOne(c *fiber.Ctx) error
	FindMany(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

type clientHandler struct {
	service Service
}

func NewClientHandler(service Service) Handler {
	return &clientHandler{service: service}
}

var SetHandler = wire.NewSet(NewClientHandler)

func (handler clientHandler) CreateOne(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)

	client := new(Client)
	if err := c.BodyParser(client); err != nil {
		return errors.HandleParseError(c, err)
	}

	result, err := handler.service.WithTx(tx).CreateOne(client)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (handler clientHandler) FindMany(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	result, err := handler.service.WithTx(tx).FindMany()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler clientHandler) DeleteOne(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id, parseErr := strconv.Atoi(c.Params("id"))
	if parseErr != nil {
		return errors.HandleParseError(c, parseErr)
	}

	result, err := handler.service.WithTx(tx).DeleteOne(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
