package slackbot

import (
	"aesir/src/channels"
	"aesir/src/common"
	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SlackbotSuit struct {
	suite.Suite
	config         *common.Config
	channelService *channels.MockService
	service        Service
	userId         string
	teamId         string
	channelId      string
	log            *logrus.Logger
}

func (suite *SlackbotSuit) SetupSuite() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	suite.config = common.NewConfig()
	suite.channelService = channels.NewMockService(suite.T())
	suite.service = NewSlackService(suite.config, suite.channelService)
	whoAmI, initializeErr := suite.service.WhoAmI()
	suite.userId = whoAmI.UserID
	suite.teamId = whoAmI.TeamID
	channelsData, initializeErr := suite.service.FindChannels()
	suite.channelId = channelsData[0].ID
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	suite.log = logger
	if initializeErr != nil {
		panic("initialize error")
	}
	println("initialized instance")
}

//	func (suite *SlackbotSuit) SetupTest() {
//		suite.log.Print("SetupTest")
//	}

func (suite *SlackbotSuit) BeforeTest(suiteName, testName string) {
	suite.channelService.ExpectedCalls = nil
	suite.channelService.Calls = nil
	suite.service = NewSlackService(suite.config, suite.channelService)
}

//func (suite *SlackbotSuit) AfterTest(suiteName, testName string) {
//}

//func (suite *SlackbotSuit) TearDownSuite() {
//	suite.log.Print("TearDownSuite")
//}

//func (suite *SlackbotSuit) TearDownTest() {
//	suite.log.Print("TearDownTest")
//}

func (suite *SlackbotSuit) TestWhoAmI() {
	data, err := suite.service.WhoAmI()
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &WhoAmI{}, data)
}

func (suite *SlackbotSuit) TestFindTeam() {
	data, err := suite.service.FindTeam()
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &slack.TeamInfo{}, data)
}

func (suite *SlackbotSuit) TestFindChannels() {
	data, err := suite.service.FindChannels()
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), []slack.Channel{}, data)
}

func (suite *SlackbotSuit) TestFindJoinedChannels() {
	data, err := suite.service.FindJoinedChannels()
	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), []slack.Channel{}, data)
}

//func (suite *SlackbotSuit) TestFindChannel() {
//	data, err := suite.service.FindChannel("test")
//	assert.Nil(suite.T(), err)
//	assert.IsType(suite.T(), &slack.Channel{}, data)
//}

//func (suite *SlackbotSuit) TestFindLatestChannelMessage() {
//	data, err := suite.service.FindLatestChannelMessage("test")
//	assert.Nil(suite.T(), err)
//	assert.IsType(suite.T(), &slack.Channel{}, data)
//}

//func (suite *SlackbotSuit) TestFindTeamUsers() {
//
//}

func (suite *SlackbotSuit) TestEventMux() {
	var memberJoinedEvt *slackevents.MemberJoinedChannelEvent
	gofakeit.Struct(&memberJoinedEvt)
	memberJoinedEvt.User = suite.userId
	memberJoinedEvt.Channel = suite.channelId
	memberJoinInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: memberJoinedEvt,
	}
	suite.channelService.On("FindOneBySlackId", mock.Anything, mock.Anything).Return(new(channels.Channel), nil)
	suite.channelService.On("UpdateOneBySlackId", mock.Anything, mock.Anything).Return(new(channels.Channel), nil)

	err1 := suite.service.EventMux(*memberJoinInnerEvt)

	messageEvt := new(slackevents.MessageEvent)
	msgInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: messageEvt,
	}
	err2 := suite.service.EventMux(*msgInnerEvt)

	assert.Nil(suite.T(), err1)
	suite.channelService.AssertNumberOfCalls(suite.T(), "FindOneBySlackId", 1)
	assert.Nil(suite.T(), err2)
}

func (suite *SlackbotSuit) TestEvt() {
	var memberJoinedEvt *slackevents.MemberJoinedChannelEvent
	gofakeit.Struct(&memberJoinedEvt)
	memberJoinedEvt.User = suite.userId
	memberJoinedEvt.Channel = suite.channelId
	memberJoinInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: memberJoinedEvt,
	}
	suite.channelService.On("FindOneBySlackId", mock.Anything, mock.Anything).Return(new(channels.Channel), nil)
	suite.channelService.On("UpdateOneBySlackId", mock.Anything, mock.Anything).Return(new(channels.Channel), nil)

	err1 := suite.service.EventMux(*memberJoinInnerEvt)

	messageEvt := new(slackevents.MessageEvent)
	msgInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: messageEvt,
	}
	err2 := suite.service.EventMux(*msgInnerEvt)

	assert.Nil(suite.T(), err1)
	suite.channelService.AssertNumberOfCalls(suite.T(), "FindOneBySlackId", 1)
	assert.Nil(suite.T(), err2)
}

func TestSlackbotSuite(t *testing.T) {
	suite.Run(t, new(SlackbotSuit))
}