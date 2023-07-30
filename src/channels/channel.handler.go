package channels

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Handler interface {
	FineMany(c *fiber.Ctx) error
	FindManyWithMessage(c *fiber.Ctx) error
	FindOneBySlackId(c *fiber.Ctx) error
}

type channelHandler struct {
	service Service
}

func NewChannelHandler(service Service) Handler {
	return &channelHandler{service}
}

var SetHandler = wire.NewSet(NewChannelHandler)

func (handler channelHandler) FineMany(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)

	result, err := handler.service.WithTx(tx).FindMany()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler channelHandler) FindManyWithMessage(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)

	result, err := handler.service.WithTx(tx).FindManyWithMessage()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler channelHandler) FindOneBySlackId(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id := c.Params("slackId")

	result, err := handler.service.WithTx(tx).FindOneBySlackId(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
