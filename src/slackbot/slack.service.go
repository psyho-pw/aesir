package slackbot

import (
	"aesir/src/common"
	"aesir/src/common/errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"log"
	"os"
)

type SlackService interface {
	EventMux(innerEvent slackevents.EventsAPIInnerEvent) error
	WhoAmI() (*WhoAmI, error)
	FindTeam() (*slack.TeamInfo, error)
	FindChannels(teamId string) ([]slack.Channel, error)
	FindChannel(channelId string) (*slack.Channel, error)
	FindLatestChannelMessage(channelId string) (*slack.Message, error)
	FindTeamUsers(teamId string) ([]slack.User, error)
	WithTx(tx *gorm.DB) SlackService
}

type slackService struct {
	api *slack.Client
}

func NewSlackService(config *common.Config) SlackService {
	api := slack.New(
		config.Slack.BotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(config.Slack.AppToken),
	)

	return &slackService{api}
}

var SetService = wire.NewSet(NewSlackService)

func (service *slackService) EventMux(innerEvent slackevents.EventsAPIInnerEvent) error {
	spew.Dump(innerEvent)

	switch evt := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		return service.handleMessageEvent(evt)
	case *slackevents.MemberJoinedChannelEvent:
		return service.memberJoinHandler(evt)
	default:
		return errors.New(fiber.StatusNoContent, "no matching event from incoming type")
	}
}

func (service *slackService) memberJoinHandler(event *slackevents.MemberJoinedChannelEvent) error {
	//TODO 봇이 채널 조인 시 채널 정보 취득하여 저장
	self, slackErr := service.WhoAmI()
	if slackErr != nil {
		return slackErr
	}
	if event.User != self.UserID {
		return nil
	}

	//TODO find channel from db & create channel
	channel, slackChannelErr := service.FindChannel(event.Channel)
	if slackChannelErr != nil {
		return slackChannelErr
	}

	spew.Dump(channel)
	return nil
}

func (service *slackService) handleMessageEvent(event *slackevents.MessageEvent) error {
	//TODO 메세지 이벤트 발생 시 사내 인원일 경우 timestamp 최신화
	if event.BotID != "" || event.ThreadTimeStamp != "" {
		return nil
	}
	_, _, err := service.api.PostMessage(event.Channel, slack.MsgOptionText("acknowledged", false))
	if err != nil {
		return err
	}
	return nil
}

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
	team, err := service.api.GetTeamInfo()
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (service *slackService) FindChannels(teamId string) ([]slack.Channel, error) {
	channels, _, err := service.api.GetConversations(
		&slack.GetConversationsParameters{
			ExcludeArchived: true,
			Limit:           1000,
			Types:           []string{"public_channel", "private_channel"},
			TeamID:          teamId,
		},
	)
	if err != nil {
		logrus.Errorf("%+v", err)
		return nil, err
	}

	return channels, nil
}

func (service *slackService) FindChannel(channelId string) (*slack.Channel, error) {
	channel, err := service.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	return channel, nil
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

	messages := getConversationHistoryResponse.Messages
	if messages == nil {
		return nil, errors.New(fiber.StatusNotFound, "latest message not found")
	}

	return &messages[0], nil
}

func (service *slackService) FindTeamUsers(teamId string) ([]slack.User, error) {
	users, err := service.api.GetUsers(slack.GetUsersOptionTeamID(teamId))
	if err != nil {
		return nil, err
	}

	pred := func(i slack.User) bool {
		return i.IsBot == false &&
			i.IsRestricted == false &&
			i.IsUltraRestricted == false &&
			i.Deleted == false
	}

	return funk.Filter(users, pred).([]slack.User), nil
}

func (service *slackService) WithTx(tx *gorm.DB) SlackService {
	return service
}
