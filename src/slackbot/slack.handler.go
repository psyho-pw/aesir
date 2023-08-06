package slackbot

import (
	_const "aesir/src/common/const"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gorm.io/gorm"
)

type Handler interface {
	EventMux(c *fiber.Ctx) error
	CommandMux(c *fiber.Ctx) error
	WhoAmI(c *fiber.Ctx) error
	FindTeam(c *fiber.Ctx) error
	FindChannels(c *fiber.Ctx) error
	FindChannelById(c *fiber.Ctx) error
	FindLatestChannelMessage(c *fiber.Ctx) error
	FindTeamUsers(c *fiber.Ctx) error
}

type slackHandler struct {
	service Service
}

func NewSlackHandler(service Service) Handler {
	return &slackHandler{service: service}
}

func (handler slackHandler) EventMux(c *fiber.Ctx) error {
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
		evtErr := handler.service.WithTx(tx).EventMux(&innerEvent)
		if evtErr != nil {
			logrus.Errorf("%+v", evtErr)
			return evtErr
		}

		return c.SendStatus(fiber.StatusOK)

	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (handler slackHandler) CommandMux(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	commandType := c.Params("commandType")

	httpRequest, convertErr := adaptor.ConvertRequest(c, false)
	if convertErr != nil {
		return convertErr
	}

	command, parseEvtErr := slack.SlashCommandParse(httpRequest)
	if parseEvtErr != nil {
		logrus.Errorf("%+v", parseEvtErr)
		return c.Status(fiber.StatusInternalServerError).SendString(parseEvtErr.Error())
	}
	logrus.Infof("%s command triggered", command.Command)

	spew.Dump(command)

	switch commandType {
	case _const.CommandTypeManager:
		logrus.Debug("manager")
		err := handler.service.WithTx(tx).ManagerCommand(command)
		if err != nil {
			return nil
		}
		break
	case _const.CommandTypeThreshold:
		logrus.Debug("threshold")
		break
	default:
		logrus.Errorf("no matching command exists")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (handler slackHandler) WhoAmI(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	result, err := handler.service.WithTx(tx).WhoAmI()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
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
	result, err := handler.service.WithTx(tx).FindChannels()
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
