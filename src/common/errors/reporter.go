package errors

import (
	_const "aesir/src/common/const"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func getCredentials(webhookUrl string) discordgo.Webhook {
	credentialsResponse, credentialsErr := http.Get(webhookUrl)
	defer func(response *http.Response) {
		if r := recover(); r != nil {
			logrus.Errorf("Reporter fatal error")
			logrus.Errorf("%+v", r)
		}
	}(credentialsResponse)

	if credentialsErr != nil {
		panic(credentialsErr)
	}

	data, readErr := io.ReadAll(credentialsResponse.Body)
	if readErr != nil {
		panic(readErr)
	}

	var credentials discordgo.Webhook
	if unmarshalErr := json.Unmarshal(data, &credentials); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	return credentials
}

func formatMessage(exception *Error) *discordgo.WebhookParams {
	logrus.Errorf("%+v", exception)
	messageField := &discordgo.MessageEmbedField{
		Name:  "Message",
		Value: exception.Message,
	}

	stackField := &discordgo.MessageEmbedField{
		Name:  "Stack",
		Value: exception.Stack,
	}

	embed := &discordgo.MessageEmbed{
		Title: func(caller *string) string {
			if caller == nil || *caller == _const.Unknown {
				return _const.UnhandledException
			}

			return *caller
		}(&exception.Caller),
		Color:  16711680,
		Fields: []*discordgo.MessageEmbedField{messageField, stackField},
	}

	params := &discordgo.WebhookParams{Embeds: []*discordgo.MessageEmbed{embed}}

	return params
}

func Report(webhookUrl string, exception *Error) error {
	credentials := getCredentials(webhookUrl)

	client, err := discordgo.New("Bot " + credentials.Token)
	if err != nil {
		fmt.Println("error creating Discord session")
		panic(err)
	}

	params := formatMessage(exception)
	_, sendErr := client.WebhookExecute(credentials.ID, credentials.Token, true, params)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
