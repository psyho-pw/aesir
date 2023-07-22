package dto

import "github.com/slack-go/slack"

type Event struct {
	Token     string      `json:"token,omitempty"`
	Challenge string      `json:"challenge,omitempty"`
	Type      string      `json:"type,omitempty"`
	Event     slack.Event `json:"event,omitempty"`
}
