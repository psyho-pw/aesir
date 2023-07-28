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

type MySuite struct {
	suite.Suite
	service Service
}

func (suite *MySuite) SetupTest() {
	//suite.Data = []string{"one", "two", "three"}
	err := godotenv.Load("../../.env")
	if err != nil {
		spew.Dump(err)
		panic("Error loading .env file")
	}
	slackService := NewSlackService(common.NewConfig(), channels.NewMockService(suite.T()))
	suite.service = slackService
	println("initialized service instance")
}

func (suite *MySuite) BeforeTest(suiteName, testName string) {
	// Run some code before each test
}

func (suite *MySuite) TestExample() {
	data, err := suite.service.WhoAmI()
	spew.Dump(data)
	assert.Equal(suite.T(), nil, err)
}

func (suite *MySuite) TestWhoAmI() {

	config := common.NewConfig()
	spew.Dump(config)
	//slackService := NewSlackService(common.NewConfig(), channels.NewMockService(t))
	//info, err := slackService.WhoAmI()
	assert.Equal(suite.T(), false, false)
	//
	//spew.Dump(info)
	// Test something using the Data fixture here
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
