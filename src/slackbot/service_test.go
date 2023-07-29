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
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SlackbotSuit struct {
	suite.Suite
	service Service
	log     *logrus.Logger
}

func (suite *SlackbotSuit) SetupSuite() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	service := NewSlackService(common.NewConfig(), channels.NewMockService(suite.T()))
	suite.service = service
	println("initialized slackService instance")

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	suite.log = logger
}

//func (suite *SlackbotSuit) SetupTest() {
//	suite.log.Print("SetupTest")
//}
//
//func (suite *SlackbotSuit) BeforeTest(suiteName, testName string) {
//	suite.log.Print("BeforeTest")
//}
//
//func (suite *SlackbotSuit) AfterTest(suiteName, testName string) {
//	suite.log.Print("AfterTest")
//}
//
//func (suite *SlackbotSuit) TearDownSuite() {
//	suite.log.Print("TearDownSuite")
//}
//
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
//
//func (suite *SlackbotSuit) TestFindLatestChannelMessage() {
//	data, err := suite.service.FindLatestChannelMessage("test")
//	assert.Nil(suite.T(), err)
//	assert.IsType(suite.T(), &slack.Channel{}, data)
//}

func (suite *SlackbotSuit) TestFindTeamUsers() {

}

func (suite *SlackbotSuit) TestEventMux() {
	var memberJoinedEvt *slackevents.MemberJoinedChannelEvent
	gofakeit.Struct(&memberJoinedEvt)
	memberJoinInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: memberJoinedEvt,
	}
	err1 := suite.service.EventMux(*memberJoinInnerEvt)

	messageEvt := new(slackevents.MessageEvent)
	msgInnerEvt := &slackevents.EventsAPIInnerEvent{
		Type: gofakeit.Word(),
		Data: messageEvt,
	}
	err2 := suite.service.EventMux(*msgInnerEvt)

	assert.Nil(suite.T(), err1)
	assert.Nil(suite.T(), err2)
}

func TestSlackbotSuite(t *testing.T) {
	suite.Run(t, new(SlackbotSuit))
}
