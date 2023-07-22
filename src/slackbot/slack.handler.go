package slackbot

import (
	"aesir/src/slackbot/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

type SlackHandler interface {
	Event(c *fiber.Ctx) error
	FindTeam(c *fiber.Ctx) error
	FindChannels(c *fiber.Ctx) error
	FindChannelById(c *fiber.Ctx) error
	FindLatestChannelMessage(c *fiber.Ctx) error
	FindTeamUsers(c *fiber.Ctx) error
}

type slackHandler struct {
	service SlackService
}

func NewSlackHandler(service SlackService) SlackHandler {
	return &slackHandler{service: service}
}

func (handler slackHandler) Event(c *fiber.Ctx) error {
	payload := new(slack.Event)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if payload.Type == "url_verification" {
		challenge := new(dto.Event)
		if err := c.BodyParser(challenge); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(challenge)
	}

	eventType := slack.EventMapping[payload.Type]
	if err := c.BodyParser(eventType); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	logrus.Debugf("%+v", payload)
	logrus.Debugf("%+v", eventType)

	return c.Status(fiber.StatusOK).JSON(payload)
}

func (handler slackHandler) FindTeam(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	result, err := handler.service.WithTx(tx).FindTeam()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler slackHandler) FindChannels(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id := c.Params("teamId")
	result, err := handler.service.WithTx(tx).FindChannels(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler slackHandler) FindChannelById(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id := c.Params("channelId")
	result, err := handler.service.WithTx(tx).FindChannel(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler slackHandler) FindLatestChannelMessage(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id := c.Params("channelId")
	result, err := handler.service.WithTx(tx).FindLatestChannelMessage(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (handler slackHandler) FindTeamUsers(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	id := c.Params("teamId")
	result, err := handler.service.WithTx(tx).FindTeamUsers(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

var SetHandler = wire.NewSet(NewSlackHandler)
