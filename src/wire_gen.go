// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package src

import (
	"fiber/src/common"
	"fiber/src/common/database"
	"fiber/src/slackbot"
	"fiber/src/users"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from wire.go:

func New() (*fiber.App, error) {
	config := common.NewConfig()
	db := database.NewDB(config)
	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(userRepository)
	userHandler := users.NewUserHandler(userService)
	slackService := slackbot.NewSlackService(config)
	slackHandler := slackbot.NewSlackHandler(slackService)
	app := NewApp(config, db, userHandler, slackHandler)
	return app, nil
}

// wire.go:

var WireSet = wire.NewSet(AppSet)
