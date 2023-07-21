package cron

import (
	"aesir/src/common"
	"aesir/src/slackbot"
	"aesir/src/users"
	"github.com/google/wire"
	"sync"
)

var svcOnce sync.Once

type CronService interface {
	Start() error
}

type cronService struct {
	config       *common.Config
	slackService slackbot.SlackService
	userService  users.UserService
}

func NewCronService(config *common.Config, slackService slackbot.SlackService, userService users.UserService) CronService {
	return &cronService{
		config:       config,
		slackService: slackService,
		userService:  userService,
	}
}

var SetService = wire.NewSet(NewCronService)

func (service *cronService) Start() error {
	teamUsers, err := service.slackService.FindTeamUsers(service.config.Slack.TeamId)
	if err != nil {
		return err
	}

	for _, user := range teamUsers {
		user, err := service.userService.FindOneBySlackId(user.ID)
		if err != nil {
			return err
		}

		if user == nil {
			newUser := new(users.User)
			newUser.Email = user.Email
			_, err := service.userService.CreateOne(newUser)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
