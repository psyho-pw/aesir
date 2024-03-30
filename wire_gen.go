// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"aesir/src"
	"aesir/src/channels"
	"aesir/src/common"
	"aesir/src/common/database"
	"aesir/src/google"
	"aesir/src/messages"
	"aesir/src/slackbot"
	"aesir/src/users"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from wire.go:

func New() (*fiber.App, error) {
	config := common.NewConfig()
	db := database.NewDB(config)
	repository := users.NewUserRepository(db)
	service := users.NewUserService(repository)
	handler := users.NewUserHandler(service)
	channelsRepository := channels.NewChannelRepository(db)
	channelsService := channels.NewChannelService(channelsRepository)
	channelsHandler := channels.NewChannelHandler(channelsService)
	messagesRepository := messages.NewMessageRepository(db)
	messagesService := messages.NewMessageService(messagesRepository)
	messagesHandler := messages.NewMessageHandler(messagesService)
	googleService := google.NewGoogleService(config)
	slackbotService := slackbot.NewSlackService(config, service, channelsService, messagesService, googleService)
	slackbotHandler := slackbot.NewSlackHandler(slackbotService)
	app := src.NewApp(config, db, handler, channelsHandler, messagesHandler, slackbotHandler, googleService)
	return app, nil
}

// wire.go:

var WireSet = wire.NewSet(src.AppSet)
