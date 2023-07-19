package slackbot

import (
	"fiber/src/common"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

type SlackService interface {
	GetChannels() ([]slack.Channel, error)
	GetChannel(channelId string) (*slack.Channel, error)
	WithTx(tx *gorm.DB) SlackService
}

type slackService struct {
	api    *slack.Client
	client *socketmode.Client
}

func socketEventListener(client *socketmode.Client) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			logrus.Debug("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			logrus.Debug("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			logrus.Debug("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				logrus.Debugf("Ignored %+v\n", evt)

				continue
			}

			logrus.Debugf("Event received: %+v\n", eventsAPIEvent)

			client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
					if err != nil {
						logrus.Debugf("failed posting message: %v", err)
					}
				case *slackevents.MemberJoinedChannelEvent:
					logrus.Debugf("user %q joined to channel %q", ev.User, ev.Channel)
				}
			default:
				client.Debugf("unsupported Events API event received")
			}
		case socketmode.EventTypeInteractive:
			callback, ok := evt.Data.(slack.InteractionCallback)
			if !ok {
				logrus.Debugf("Ignored %+v\n", evt)

				continue
			}

			logrus.Debugf("Interaction received: %+v\n", callback)

			var payload interface{}

			switch callback.Type {
			case slack.InteractionTypeBlockActions:
				// See https://api.slack.com/apis/connections/socket-implement#button

				client.Debugf("button clicked!")
			case slack.InteractionTypeShortcut:
			case slack.InteractionTypeViewSubmission:
				// See https://api.slack.com/apis/connections/socket-implement#modal
			case slack.InteractionTypeDialogSubmission:
			default:

			}

			client.Ack(*evt.Request, payload)
		case socketmode.EventTypeSlashCommand:
			cmd, ok := evt.Data.(slack.SlashCommand)
			if !ok {
				logrus.Debugf("Ignored %+v\n", evt)

				continue
			}

			client.Debugf("Slash command received: %+v", cmd)

			payload := map[string]interface{}{
				"blocks": []slack.Block{
					slack.NewSectionBlock(
						&slack.TextBlockObject{
							Type: slack.MarkdownType,
							Text: "foo",
						},
						nil,
						slack.NewAccessory(
							slack.NewButtonBlockElement(
								"",
								"somevalue",
								&slack.TextBlockObject{
									Type: slack.PlainTextType,
									Text: "bar",
								},
							),
						),
					),
				},
			}

			client.Ack(*evt.Request, payload)
		default:
			logrus.Debugf("Unexpected event type received: %s\n", evt.Type)
		}
	}
}

func NewSlackService(config *common.Config) SlackService {
	appToken := config.Slack.AppToken
	botToken := config.Slack.BotToken

	if appToken == "" {
		logrus.Error("Missing slack app token")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		logrus.Error("app token must have the prefix \"xapp-\"")
	}

	if botToken == "" {
		logrus.Error("Missing slack bot token")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		logrus.Error("bot token must have the prefix \"xoxb-\"")
	}

	api := slack.New(
		botToken, slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go socketEventListener(client)

	go func() {
		err := client.Run()
		if err != nil {
			logrus.Fatalf("%+v", err)
			os.Exit(1)
		}
	}()

	return &slackService{
		api:    api,
		client: client,
	}

}

var SetService = wire.NewSet(NewSlackService)

func (service *slackService) GetChannels() ([]slack.Channel, error) {
	channels, _, err := service.client.GetConversations(
		&slack.GetConversationsParameters{
			ExcludeArchived: true,
			Limit:           1000,
			Types:           nil,
		},
	)
	if err != nil {
		logrus.Errorf("%s\n", err)
		return nil, err
	}

	return channels, nil
}

func (service *slackService) GetChannel(channelId string) (*slack.Channel, error) {
	channel, err := service.client.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelId,
	})
	if err != nil {
		return nil, err
	}

	return channel, err
}

func (service *slackService) WithTx(tx *gorm.DB) SlackService {
	return service
}
