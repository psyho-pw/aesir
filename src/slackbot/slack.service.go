package slackbot

import (
	"fiber/src/common"
	"github.com/google/wire"
	"github.com/slack-go/slack"
)

type SlackService interface {
	GetChannels() error
}

type slackService struct {
	slackBotInstance *slack.Client
}

func NewSlackService(config *common.Config) SlackService {
	println(":::::::::::::::::::::::::::::::")
	return &slackService{slackBotInstance: slack.New(config.Slack.AppToken)}
}

var SetService = wire.NewSet(NewSlackService)

func (service *slackService) GetChannels() error {
	return nil
}
