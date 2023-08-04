package messages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Handler interface {
	FindMany(c *fiber.Ctx) error
}

type messageHandler struct {
	service Service
}

func NewMessageHandler(service Service) Handler {
	return &messageHandler{service}
}

var SetHandler = wire.NewSet(NewMessageHandler)

func (handler *messageHandler) FindMany(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)

	result, err := handler.service.WithTx(tx).FindMany()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
