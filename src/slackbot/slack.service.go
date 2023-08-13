package slackbot

import (
	"aesir/src/channels"
	"aesir/src/common"
	_const "aesir/src/common/const"
	"aesir/src/common/errors"
	"aesir/src/common/utils"
	"aesir/src/messages"
	"aesir/src/users"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"strconv"
)

type Service interface {
	EventMux(innerEvent *slackevents.EventsAPIInnerEvent) error
	OnManagerCommand(command slack.SlashCommand) error
	OnThresholdCommand(command slack.SlashCommand) error
	OnInteractionTypeManagerSelect(selectedOptions *[]slack.OptionBlockObject) error
	OnInteractionTypeThresholdSelect(selectedOption *slack.OptionBlockObject) error
	WhoAmI() (*WhoAmI, error)
	FindTeam() (*slack.TeamInfo, error)
	FindChannels() ([]slack.Channel, error)
	FindJoinedChannels() ([]slack.Channel, error)
	FindChannel(channelId string) (*slack.Channel, error)
	FindLatestChannelMessage(channelId string) (*slack.Message, error)
	FindTeamUsers(teamId string) ([]slack.User, error)
	WithTx(tx *gorm.DB) Service
}

type slackService struct {
	api            *slack.Client
	userService    users.Service
	channelService channels.Service
	messageService messages.Service
}

func NewSlackService(config *common.Config, userService users.Service, channelService channels.Service, messageService messages.Service) Service {
	api := slack.New(
		config.Slack.BotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(config.Slack.AppToken),
	)

	return &slackService{api, userService, channelService, messageService}
}

var SetService = wire.NewSet(NewSlackService)

/*
********** Event related services
 */

func (service *slackService) EventMux(innerEvent *slackevents.EventsAPIInnerEvent) error {
	switch evt := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		return service.handleMessageEvent(evt)
	case *slackevents.MemberJoinedChannelEvent:
		return service.handleMemberJoinEvent(evt)
	default:
		logrus.Debugf("no matching event from incoming type")
		return nil
	}
}

func (service *slackService) handleMemberJoinEvent(event *slackevents.MemberJoinedChannelEvent) error {
	self, slackErr := service.WhoAmI()
	if slackErr != nil {
		return slackErr
	}
	if event.User != self.UserID {
		return nil
	}

	channel, slackChannelErr := service.FindChannel(event.Channel)
	if slackChannelErr != nil {
		return slackChannelErr
	}

	persistentChannel, err := service.channelService.FindOneBySlackId(channel.ID)
	if err != nil {
		return err
	}
	if persistentChannel != nil {
		newChannel := new(channels.Channel)
		newChannel.SlackId = channel.ID
		newChannel.Name = channel.Name
		newChannel.Creator = channel.Creator
		newChannel.IsPrivate = channel.IsPrivate
		newChannel.IsArchived = channel.IsArchived

		isSame := reflect.DeepEqual(persistentChannel, newChannel)
		if isSame == true {
			return nil
		}

		_, updateErr := service.channelService.UpdateOneBySlackId(channel.ID, *newChannel)
		if updateErr != nil {
			return updateErr
		}

		return nil
	}

	newChannel := new(channels.Channel)
	newChannel.SlackId = channel.ID
	newChannel.Name = channel.Name
	newChannel.Creator = channel.Creator
	newChannel.IsPrivate = channel.IsPrivate
	newChannel.IsArchived = channel.IsArchived

	_, creationErr := service.channelService.Create(*newChannel)
	if creationErr != nil {
		return creationErr
	}

	return nil
}

// 메세지 이벤트 발생 시 화자가 사내 인원일 경우 message timestamp 최신화
// 사내 인원이 아닐 경우 pass
func (service *slackService) handleMessageEvent(event *slackevents.MessageEvent) error {
	if event.BotID != "" || event.ThreadTimeStamp != "" {
		return nil
	}

	user, fetchUserErr := service.userService.FindOneBySlackId(event.User)
	if fetchUserErr != nil {
		return fetchUserErr
	}

	if user == nil {
		logrus.Debugf("user %s is not a member of workspace", event.User)
	}

	channel, fetchChannelErr := service.channelService.FindOneBySlackId(event.Channel)
	if fetchChannelErr != nil {
		return fetchChannelErr
	}

	logrus.Infof("<%s>[%s - %s]: %s (%s)", channel.Name, user.Name, event.User, event.Text, event.TimeStamp)

	if channel.Message == nil {
		message := new(messages.Message)
		message.ChannelId = channel.ID
		channel.Message = message
	}

	var parseErr error
	channel.Message.Content = event.Text
	channel.Message.Timestamp, parseErr = strconv.ParseFloat(event.TimeStamp, 64)
	if parseErr != nil {
		return parseErr
	}

	_, channelUpdateErr := service.channelService.UpdateOneBySlackId(event.Channel, *channel)
	if channelUpdateErr != nil {
		return channelUpdateErr
	}

	return nil
}

/*
********** Command related services
 */

func (service *slackService) makeModalElements(selectElement interface{}) (*slack.TextBlockObject, *slack.TextBlockObject, *slack.Blocks, error) {
	ref := reflect.ValueOf(selectElement).Elem()
	typeofRef := ref.Type()

	titleText := slack.NewTextBlockObject("plain_text", "Aesir", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Designate a person to receive contact", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	selectLabel := slack.NewTextBlockObject("plain_text", "Threshold", false, false)

	var selectSection *slack.SectionBlock
	switch typeofRef.Name() {
	case "SelectBlockElement":
		selectSection = slack.NewSectionBlock(selectLabel, nil, slack.NewAccessory(selectElement.(*slack.SelectBlockElement)))
		break
	case "MultiSelectBlockElement":
		selectSection = slack.NewSectionBlock(selectLabel, nil, slack.NewAccessory(selectElement.(*slack.MultiSelectBlockElement)))
		break
	default:
		return nil, nil, nil, errors.New(fiber.StatusInternalServerError, "no matching type")
	}

	blocks := slack.Blocks{BlockSet: []slack.Block{headerSection, selectSection}}

	return titleText, closeText, &blocks, nil
}

func (service *slackService) OnManagerCommand(command slack.SlashCommand) error {
	var options []*slack.OptionBlockObject

	usersData, fetchUserErr := service.userService.FindMany()
	if fetchUserErr != nil {
		return fetchUserErr
	}

	var managers []users.User

	for _, user := range usersData {
		if user.IsManager == true {
			managers = append(managers, user)
		}

		option := slack.NewOptionBlockObject(
			strconv.Itoa(int(user.ID)),
			slack.NewTextBlockObject("plain_text", user.Name, false, false),
			nil,
		)
		options = append(options, option)
	}

	selectPlaceholder := slack.NewTextBlockObject("plain_text", "select..", false, false)
	multiSelectElement := slack.NewOptionsMultiSelectBlockElement(
		"multi_static_select",
		selectPlaceholder,
		_const.InteractionTypeManagerSelect,
		options...,
	)

	//set max selected item count
	maxSelectedItems := _const.MaxSelectedItems
	multiSelectElement.MaxSelectedItems = &maxSelectedItems

	// set already selected managers as initial options
	if len(managers) > 0 {
		var initialOptions []*slack.OptionBlockObject
		for _, manager := range managers {
			option := slack.NewOptionBlockObject(
				strconv.Itoa(int(manager.ID)),
				slack.NewTextBlockObject("plain_text", manager.Name, false, false),
				nil,
			)
			initialOptions = append(initialOptions, option)
		}

		multiSelectElement.InitialOptions = initialOptions
	}

	titleText, closeText, blocks, makeModalErr := service.makeModalElements(multiSelectElement)
	if makeModalErr != nil {
		return makeModalErr
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Blocks = *blocks

	_, err := service.api.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return err
	}

	return nil
}

func (service *slackService) OnThresholdCommand(command slack.SlashCommand) error {
	predicate := func(i int) *slack.OptionBlockObject {
		option := slack.NewOptionBlockObject(
			strconv.Itoa(i),
			slack.NewTextBlockObject("plain_text", fmt.Sprintf("%.2d day (s)", i), false, false),
			nil,
		)

		return option
	}
	options := funk.Map(utils.MakeRange(1, 10), predicate).([]*slack.OptionBlockObject)

	selectPlaceholder := slack.NewTextBlockObject("plain_text", "select..", false, false)
	selectElement := slack.NewOptionsSelectBlockElement(
		"static_select",
		selectPlaceholder,
		_const.InteractionTypeThresholdSelect,
		options...,
	)

	//TODO set already selected threshold as initial option
	//if len(managers) > 0 {
	//	var initialOptions []*slack.OptionBlockObject
	//	for _, manager := range managers {
	//		option := slack.NewOptionBlockObject(
	//			strconv.Itoa(int(manager.ID)),
	//			slack.NewTextBlockObject("plain_text", manager.Name, false, false),
	//			nil,
	//		)
	//		initialOptions = append(initialOptions, option)
	//	}
	//
	//	multiSelectElement.InitialOptions = initialOptions
	//}

	titleText, closeText, blocks, makeModalErr := service.makeModalElements(selectElement)
	if makeModalErr != nil {
		return makeModalErr
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Blocks = *blocks

	_, err := service.api.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return err
	}

	return nil
}

/*
********** Interaction related services
 */

func (service *slackService) OnInteractionTypeManagerSelect(selectedOptions *[]slack.OptionBlockObject) error {
	predicate := func(i slack.OptionBlockObject) int {
		id, _ := strconv.Atoi(i.Value)
		return id
	}

	userIds := funk.Map(*selectedOptions, predicate)
	err := service.userService.UpdateManagers(userIds.([]int))
	if err != nil {
		return err
	}

	return nil
}

func (service *slackService) OnInteractionTypeThresholdSelect(selectedOption *slack.OptionBlockObject) error {
	value, _ := strconv.Atoi(selectedOption.Value)
	err := service.channelService.UpdateThreshold(value)
	if err != nil {
		return err
	}

	return nil
}

/*
********** Slack API related services
 */

type WhoAmI struct {
	TeamID string `json:"teamId"`
	UserID string `json:"userId"`
	BotID  string `json:"botId"`
}

func (service *slackService) WhoAmI() (*WhoAmI, error) {
	authTestResponse, err := service.api.AuthTest()
	if err != nil {
		return nil, err
	}

	whoAmI := &WhoAmI{
		TeamID: authTestResponse.TeamID,
		UserID: authTestResponse.UserID,
		BotID:  authTestResponse.BotID,
	}
	return whoAmI, nil
}

func (service *slackService) FindTeam() (*slack.TeamInfo, error) {
	return service.api.GetTeamInfo()
}

func (service *slackService) FindChannels() ([]slack.Channel, error) {
	whoAmI, authErr := service.WhoAmI()
	if authErr != nil {
		return nil, authErr
	}

	channelsData, _, err := service.api.GetConversations(
		&slack.GetConversationsParameters{
			ExcludeArchived: true,
			Limit:           1000,
			Types:           []string{"public_channel", "private_channel"},
			TeamID:          whoAmI.TeamID,
		},
	)
	if err != nil {
		return nil, err
	}

	return channelsData, nil
}

func (service *slackService) FindJoinedChannels() ([]slack.Channel, error) {
	whoAmI, authErr := service.WhoAmI()
	if authErr != nil {
		return nil, authErr
	}

	channelsData, _, err := service.api.GetConversations(
		&slack.GetConversationsParameters{
			ExcludeArchived: true,
			Limit:           1000,
			Types:           []string{"public_channel", "private_channel"},
			TeamID:          whoAmI.TeamID,
		},
	)
	if err != nil {
		return nil, err
	}

	predicate := func(i slack.Channel) bool {
		return i.IsChannel == true && i.IsMember == true
	}

	return funk.Filter(channelsData, predicate).([]slack.Channel), nil
}

func (service *slackService) FindChannel(channelId string) (*slack.Channel, error) {
	return service.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelId,
	})
}

func (service *slackService) FindLatestChannelMessage(channelId string) (*slack.Message, error) {
	getConversationHistoryResponse, err := service.api.GetConversationHistory(
		&slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Limit:     1,
		},
	)
	if err != nil {
		return nil, err
	}

	messagesResponse := getConversationHistoryResponse.Messages
	if messagesResponse == nil {
		return nil, errors.New(fiber.StatusNotFound, "latest message not found")
	}

	return &messagesResponse[0], nil
}

func (service *slackService) FindTeamUsers(teamId string) ([]slack.User, error) {
	usersResponse, err := service.api.GetUsers(slack.GetUsersOptionTeamID(teamId))
	if err != nil {
		return nil, err
	}

	predicate := func(i slack.User) bool {
		return i.IsBot == false &&
			i.IsRestricted == false &&
			i.IsUltraRestricted == false &&
			i.Deleted == false &&
			i.ID != "USLACKBOT"
	}

	return funk.Filter(usersResponse, predicate).([]slack.User), nil
}

func (service *slackService) WithTx(tx *gorm.DB) Service {
	service.userService = service.userService.WithTx(tx)
	service.channelService = service.channelService.WithTx(tx)
	service.messageService = service.messageService.WithTx(tx)

	return service
}
