package cron

import (
	"aesir/src/channels"
	"aesir/src/common"
	"aesir/src/common/errors"
	"aesir/src/slackbot"
	"aesir/src/users"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

var svcOnce sync.Once
var svc *cronService

type Service interface {
	Start() error
}

type cronService struct {
	config         *common.Config
	db             *gorm.DB
	slackService   slackbot.Service
	userService    users.Service
	channelService channels.Service
}

func New(
	config *common.Config,
	db *gorm.DB,
	slackService slackbot.Service,
	userService users.Service,
	channelService channels.Service,
) Service {
	svcOnce.Do(func() {
		svc = &cronService{
			config,
			db,
			slackService,
			userService,
			channelService,
		}
	})

	return svc
}

var SetService = wire.NewSet(New)

func (service *cronService) transactionWrapper(fn func(tx *gorm.DB) error) func() {
	return func() {
		tx := service.db.Begin()
		logrus.Info("Transaction start")

		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("%+v", r)
				//TODO error reporting
				tx.Rollback()
				logrus.Error("Transaction rollback")
			}
			logrus.Debug("Transaction end")
		}()

		err := fn(tx)
		if err != nil {
			panic(err)
		}
	}
}

func (service *cronService) userTask(tx *gorm.DB) error {
	logrus.Infof("running userTask")
	teamUsers, findUsersErr := service.slackService.WithTx(tx).FindTeamUsers(service.config.Slack.TeamId)
	if findUsersErr != nil {
		return findUsersErr
	}

	logrus.Info(len(teamUsers))

	for index, user := range teamUsers {
		usr, findUserErr := service.userService.WithTx(tx).FindOneBySlackId(user.ID)

		logrus.Infof("%+v", user)
		logrus.Debugf("%+v", usr)
		if findUserErr != nil {
			logrus.Errorf("%+v", findUserErr)
			return findUserErr
		}

		if usr == nil {
			newUser := new(users.User)
			newUser.SlackId = user.ID
			newUser.Alias = user.Name
			newUser.Name = user.RealName
			newUser.Email = user.Profile.Email
			newUser.Phone = user.Profile.Phone
			_, err := service.userService.WithTx(tx).CreateOne(newUser)
			if err != nil {
				logrus.Errorf("%+v", err)
				return err
			}
		}

		if index == 12 {
			return errors.New(fiber.StatusBadGateway, "test")
		}
	}
	return nil
}

func (service *cronService) channelTask(tx *gorm.DB) error {
	logrus.Infof("running channelTask")
	return nil
}

func (service *cronService) Start() error {
	scheduler := gocron.NewScheduler(time.Local)
	_, _ = scheduler.Every(1).Minute().Do(service.transactionWrapper(service.userTask))
	_, _ = scheduler.Every(1).Minute().Do(service.transactionWrapper(service.channelTask))

	scheduler.StartBlocking()

	return nil
}
