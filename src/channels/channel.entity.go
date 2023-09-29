package channels

import (
	"aesir/src/messages"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	SlackId    string            `json:"slackId"`
	Name       string            `json:"name"`
	Creator    string            `json:"creator"`
	IsPrivate  bool              `json:"isPrivate"`
	IsArchived bool              `json:"isArchived"`
	Message    *messages.Message `json:"message"`
	Threshold  int               `json:"threshold"`
}
