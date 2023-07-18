package slackbot

import "github.com/google/wire"

type SlackHandler interface {
}

type slackHandler struct {
	service SlackService
}

func NewSlackHandler(service SlackService) SlackHandler {
	return &slackHandler{service: service}
}

var SetHandler = wire.NewSet(NewSlackHandler)
