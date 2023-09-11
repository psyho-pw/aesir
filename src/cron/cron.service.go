package cron

import (
	"aesir/src/channels"
	"aesir/src/common"
	"aesir/src/common/utils"
	"aesir/src/messages"
	"aesir/src/slackbot"
	"aesir/src/users"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
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
	messageService messages.Service
}

func New(
	config *common.Config,
	db *gorm.DB,
	slackService slackbot.Service,
	userService users.Service,
	channelService channels.Service,
	messageService messages.Service,
) Service {
	svcOnce.Do(func() {
		svc = &cronService{
			config,
			db,
			slackService,
			userService,
			channelService,
			messageService,
		}
	})

	return svc
}

var SetService = wire.NewSet(New)

func (service *cronService) isWeekendOrHoliday() bool {
	now := time.Now()
	//if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
	//	logrus.Infof("It's weekend!")
	//	return true
	//}

	uri, uriBuildErr := service.config.OpenApi.GetUrl(now)
	if uriBuildErr != nil {
		panic(uri)
	}

	logrus.Printf("%s", uri)
	response, openApiErr := http.Get(uri)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if openApiErr != nil {
		panic(openApiErr)
	}

	data, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		panic(readErr)
	}
	var openApiResponse OpenApiResponse
	if unMarshalErr := json.Unmarshal(data, &openApiResponse); unMarshalErr != nil {
		panic(unMarshalErr)
	}

	// check holidays
	for _, item := range openApiResponse.Response.Body.Items.Item {
		itemDate, err := time.Parse("20060102", strconv.Itoa(item.LocDate))
		if err != nil {
			logrus.Errorf("date parse err")
			continue
		}

		if now.Day() == itemDate.Day() {
			logrus.Infof("It's holiday!")
			return true
		}
	}

	return false
}

func (service *cronService) transactionWrapper(fn func(tx *gorm.DB) error) func() {
	return func() {
		if service.isWeekendOrHoliday() == true {
			return
		}

		tx := service.db.Begin()
		logrus.Info("Transaction start")

		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("%+v", r)
				//TODO error reporting
				tx.Rollback()
				logrus.Error("Transaction rollback")
			}

			tx.Commit()
			logrus.Debug("Transaction end")
		}()

		err := fn(tx)
		if err != nil {
			panic(err)
		}
	}
}

func (service *cronService) userTask(tx *gorm.DB) error {
	defer utils.Timer()()
	logrus.Infof("running userTask")
	teamUsers, findUsersErr := service.slackService.WithTx(tx).FindTeamUsers(service.config.Slack.TeamId)
	if findUsersErr != nil {
		return findUsersErr
	}

	logrus.Debugf("found %d users", len(teamUsers))
	var toCreate []users.User
	for _, user := range teamUsers {
		usr, findUserErr := service.userService.WithTx(tx).FindOneBySlackId(user.ID)
		if findUserErr != nil {
			return findUserErr
		}
		if usr != nil {
			continue
		}

		newUser := new(users.User)
		newUser.SlackId = user.ID
		newUser.Alias = user.Name
		newUser.Name = user.RealName
		newUser.Email = user.Profile.Email
		newUser.Phone = user.Profile.Phone

		toCreate = append(toCreate, *newUser)
		logrus.Debugf("%+v", newUser)
	}

	if len(toCreate) == 0 {
		logrus.Debug("user target not found")
		return nil
	}

	_, insertErr := service.userService.WithTx(tx).CreateMany(toCreate)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (service *cronService) channelTask(tx *gorm.DB) error {
	defer utils.Timer()()
	logrus.Infof("running channelTask")
	joinedChannels, err := service.slackService.WithTx(tx).FindJoinedChannels()
	if err != nil {
		return err
	}

	logrus.Debugf("found %d channels", len(joinedChannels))
	var toCreate []channels.Channel
	for _, channel := range joinedChannels {
		ch, findChannelErr := service.channelService.WithTx(tx).FindOneBySlackId(channel.ID)
		if findChannelErr != nil {
			return findChannelErr
		}
		if ch != nil {
			continue
		}

		newChannel := new(channels.Channel)
		newChannel.SlackId = channel.ID
		newChannel.Name = channel.Name
		newChannel.Creator = channel.Creator
		newChannel.IsPrivate = channel.IsPrivate
		newChannel.IsArchived = channel.IsArchived

		toCreate = append(toCreate, *newChannel)
		logrus.Debugf("%+v", newChannel)
	}

	if len(toCreate) == 0 {
		logrus.Debug("channel target not found")
		return nil
	}

	_, insertErr := service.channelService.WithTx(tx).CreateMany(toCreate)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (service *cronService) notificationTask(tx *gorm.DB) error {
	defer utils.Timer()()
	logrus.Infof("running notificationTask")
	channel, findFirstErr := service.channelService.FindFirstOne()
	if findFirstErr != nil {
		return findFirstErr
	}

	targetChannels, err := service.channelService.FindManyByThreshold(channel.Threshold)
	if err != nil {
		return err
	}

	if len(targetChannels) == 0 {
		logrus.Infof("there are no channels to notified")
		return nil
	}

	predicate := func(i channels.Channel) string {
		return i.Name
	}
	channelNames := funk.Map(targetChannels, predicate).([]string)
	utils.PrettyPrint(targetChannels)
	sendDMErr := service.slackService.SendDM(channelNames)
	if sendDMErr != nil {
		return sendDMErr
	}

	idPredicate := func(i channels.Channel) uint {
		return i.ID
	}

	channelIds := funk.Map(targetChannels, idPredicate).([]int)
	service.messageService.UpdateTimestampsByChannelIds(channelIds, targetChannels[0].Threshold)

	return nil
}

func (service *cronService) Start() error {
	//if service.config.AppEnv != "production" {
	//	return nil
	//}
	scheduler := gocron.NewScheduler(time.Local)
	//_, _ = scheduler.CronWithSeconds("0 * * * * *").Do(service.transactionWrapper(service.userTask))
	//_, _ = scheduler.Every(5).Minute().Do(service.transactionWrapper(service.userTask))
	//_, _ = scheduler.Every(5).Minute().Do(service.transactionWrapper(service.channelTask))
	_, _ = scheduler.Every(5).Minute().Do(service.transactionWrapper(service.notificationTask))
	scheduler.StartAsync()

	return nil
}
