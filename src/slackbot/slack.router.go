package slackbot

import (
	"fiber/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler SlackHandler) {
	router.Get("/slack", middlewares.TxMiddleware(db), handler.FindChannels)
}
