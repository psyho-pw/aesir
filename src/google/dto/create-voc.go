package dto

import (
	"aesir/src/clients"
	"github.com/slack-go/slack"
)

type CreateVoCDto struct {
	User          *slack.User
	Client        *clients.Client
	IsStakeholder bool
	VocContent    string
}
