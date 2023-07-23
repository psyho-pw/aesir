package crons

import (
	"aesir/src/channels"
	"aesir/src/common"
	"aesir/src/slackbot"
	"aesir/src/users"
	"github.com/go-co-op/gocron"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var svcOnce sync.Once
var svc *cronService

type CronService interface {
	Start() error
}

type cronService struct {
	config         *common.Config
	slackService   slackbot.SlackService
	userService    users.UserService
	channelService channels.ChannelService
}

func NewCronService(
	config *common.Config,
	slackService slackbot.SlackService,
	userService users.UserService,
	channelService channels.ChannelService,
) CronService {
	svcOnce.Do(func() {
		svc = &cronService{config, slackService, userService, channelService}
	})

	return svc
}

var SetService = wire.NewSet(NewCronService)

func (service *cronService) channelTask() error {
	return nil
}

func (service *cronService) userTask() error {
	teamUsers, err := service.slackService.FindTeamUsers(service.config.Slack.TeamId)
	if err != nil {
		return err
	}

	println(len(teamUsers))

	for _, user := range teamUsers {
		usr, err := service.userService.FindOneBySlackId(user.ID)

		logrus.Infof("%+v", user)
		logrus.Debugf("%+v", usr)
		if err != nil {
			logrus.Errorf("%+v", err)
			return err
		}

		if usr.ID == 0 {
			newUser := new(users.User)
			newUser.SlackId = user.ID
			newUser.Alias = user.Name
			newUser.Name = user.RealName
			newUser.Email = user.Profile.Email
			newUser.Phone = user.Profile.Phone
			_, err := service.userService.CreateOne(newUser)
			if err != nil {
				logrus.Errorf("%+v", err)
				return err
			}
		}
	}
	return nil
}

func (service *cronService) Start() {
	scheduler := gocron.NewScheduler(time.Local)
	_, _ = scheduler.Every(1).Minute().Do(service.userTask())
	_, _ = scheduler.Every(1).Minute().Do(service.channelTask())

	scheduler.StartBlocking()
}
