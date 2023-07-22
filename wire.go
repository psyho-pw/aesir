//go:build wireinject
// +build wireinject

package main

import (
	"aesir/src"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(src.AppSet)

func New() (*fiber.App, error) {
	wire.Build(WireSet)
	return nil, nil
}
