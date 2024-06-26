package slack

import (
	_const "aesir/src/common/const"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gorm.io/gorm"
	"net/url"
	"strconv"
)

type Handler interface {
	EventMux(c *fiber.Ctx) error
	CommandMux(c *fiber.Ctx) error
	InteractionMux(c *fiber.Ctx) error
	WhoAmI(c *fiber.Ctx) error
	FindTeam(c *fiber.Ctx) error
	FindChannels(c *fiber.Ctx) error
	FindChannelById(c *fiber.Ctx) error
	FindLatestChannelMessage(c *fiber.Ctx) error
	FindTeamUsers(c *fiber.Ctx) error
	FindSheet(c *fiber.Ctx) error
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

	var responseStr string
	var err error
	switch commandType {
	case _const.CommandTypeManager:
		err = handler.service.WithTx(tx).OnManagerCommand(command)
		break
	case _const.CommandTypeThreshold:
		err = handler.service.WithTx(tx).OnThresholdCommand(command)
		break
	case _const.CommandTypeLeave:
		err = handler.service.WithTx(tx).OnLeaveCommand(command)
		responseStr = "was removed from channel"
		break
	case _const.CommandTypeVoC:
		err = handler.service.WithTx(tx).OnVoCCommand(command)
		break
	case _const.CommandTypeClient:
		err = handler.service.WithTx(tx).OnClientCommand(command)
	default:
		logrus.Errorf("no matching command exists")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(responseStr)
}

func (handler slackHandler) handleBlockAction(c *fiber.Ctx, message *slack.InteractionCallback, tx *gorm.DB) error {
	action := *message.ActionCallback.BlockActions[0]
	logrus.Debugf("%v", action.Value)
	switch action.ActionID {
	case _const.InteractionTypeManagerSelect:
		return handler.service.WithTx(tx).OnInteractionTypeManagerSelect(&action.SelectedOptions)

	case _const.InteractionTypeThresholdSelect:
		return handler.service.WithTx(tx).OnInteractionTypeThresholdSelect(&action.SelectedOption)

	case _const.InteractionTypeClientCreate:
		return handler.service.WithTx(tx).OnInteractionTypeClientCreate(message.View.ID, message.View.State)

	case _const.InteractionTypeClientDelete:
		clientId, err := strconv.Atoi(action.Value)
		if err != nil {
			return err
		}
		return handler.service.WithTx(tx).OnInteractionTypeClientDelete(message.View.ID, clientId)

	default:
		logrus.Errorf("no matching block action handler exists")
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func (handler slackHandler) handleViewSubmission(c *fiber.Ctx, message *slack.InteractionCallback, tx *gorm.DB) error {
	callBackId := message.View.CallbackID

	var err error
	switch callBackId {
	case _const.InteractionTypeVoCViewSubmit:
		err = handler.service.WithTx(tx).OnInteractionTypeVoCViewSubmit(&message.User, message.View.State)
		break
	default:
		logrus.Errorf("no matching view submission handler exists")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func (handler slackHandler) InteractionMux(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	jsonStr, escapeErr := url.QueryUnescape(string(c.Body())[8:])
	if escapeErr != nil {
		logrus.Errorf("[ERROR] Failed to unescape request body: %s", escapeErr)
		return escapeErr
	}

	var message *slack.InteractionCallback
	if unmarshalErr := json.Unmarshal([]byte(jsonStr), &message); unmarshalErr != nil {
		logrus.Errorf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		return unmarshalErr
	}

	//utils.PrettyPrint(message)

	switch message.Type {
	case slack.InteractionTypeBlockActions:
		return handler.handleBlockAction(c, message, tx)
	case slack.InteractionTypeViewSubmission:
		return handler.handleViewSubmission(c, message, tx)

	default:
		logrus.Errorf("no matching interaction handler exists")
		return c.SendStatus(fiber.StatusBadRequest)
	}
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

func (handler slackHandler) FindSheet(c *fiber.Ctx) error {
	tx := c.Locals("TX").(*gorm.DB)
	result, err := handler.service.WithTx(tx).FindSheet()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

var SetHandler = wire.NewSet(NewSlackHandler)
