package slackbot

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler SlackHandler) {
	router.Get("/teams", middlewares.TxMiddleware(db), handler.FindTeam)
	router.Get("/teams/:teamId/channels", middlewares.TxMiddleware(db), handler.FindChannels)
	router.Get("/teams/:teamId/channels/:channelId", middlewares.TxMiddleware(db), handler.FindChannelById)
	router.Get("/teams/:teamId/channels/:channelId/messages/latest", middlewares.TxMiddleware(db), handler.FindLatestChannelMessage)
}
