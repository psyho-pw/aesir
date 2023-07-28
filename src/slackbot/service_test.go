package slackbot

import (
	"aesir/src/channels"
	"aesir/src/common"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SlackbotSuit struct {
	suite.Suite
	service Service
}

func (suite *SlackbotSuit) SetupTest() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	slackService := NewSlackService(common.NewConfig(), channels.NewMockService(suite.T()))
	suite.service = slackService
	println("initialized service instance")
}

func (suite *SlackbotSuit) BeforeTest(suiteName, testName string) {
	// Run some code before each test
}

func (suite *SlackbotSuit) TestWhoAmI() {
	data, err := suite.service.WhoAmI()
	spew.Dump(data)
	assert.Equal(suite.T(), nil, err)
	assert.IsType(suite.T(), &WhoAmI{}, data)
}

func (suite *SlackbotSuit) TestFindTeam() {

}

func (suite *SlackbotSuit) TestFindChannels() {

}

func (suite *SlackbotSuit) TestFindJoinedChannels() {

}

func (suite *SlackbotSuit) TestFindJChannel() {

}

func (suite *SlackbotSuit) FindLatestChannelMessage() {

}

func (suite *SlackbotSuit) FindTeamUsers() {

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SlackbotSuit))
}
