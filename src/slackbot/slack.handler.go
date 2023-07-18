package slackbot

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

type SlackHandler interface {
	FindChannels(c *fiber.Ctx) error
}

type slackHandler struct {
	service SlackService
}

func NewSlackHandler(service SlackService) SlackHandler {
	return &slackHandler{service: service}
}

func (handler slackHandler) FindChannels(c *fiber.Ctx) error {
	_ = handler.service.GetChannels()
	return nil
}

var SetHandler = wire.NewSet(NewSlackHandler)
