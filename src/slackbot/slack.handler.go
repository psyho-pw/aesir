package slackbot

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
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
	tx := c.Locals("TX").(*gorm.DB)

	eventsAPIEvent, parseEvtErr := slackevents.ParseEvent(json.RawMessage(c.Body()), slackevents.OptionNoVerifyToken())
	if parseEvtErr != nil {
		logrus.Errorf("%+v", parseEvtErr)
		return c.Status(fiber.StatusInternalServerError).SendString(parseEvtErr.Error())
	}
	logrus.Infof("%s event triggered", eventsAPIEvent.Type)

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		return c.Status(fiber.StatusOK).JSON(eventsAPIEvent.Data)
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		evtErr := handler.service.WithTx(tx).HandleEvent(innerEvent)
		if evtErr != nil {
			logrus.Errorf("%+v", evtErr)
			return evtErr
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return nil
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
