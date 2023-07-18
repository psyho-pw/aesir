package slackbot

import (
	"fiber/src/common"
	"fmt"
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
	return &slackService{slackBotInstance: slack.New(config.Slack.Token)}

}

var SetService = wire.NewSet(NewSlackService)

func (service *slackService) GetChannels() error {
	groups, err := service.slackBotInstance.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}
	return nil
}
