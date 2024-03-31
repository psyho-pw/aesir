package slack

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler Handler) {
	router.Post("/events", middlewares.TxMiddleware(db), handler.EventMux)
	router.Post("/commands/:commandType", middlewares.TxMiddleware(db), handler.CommandMux)
	router.Post("/interactions", middlewares.TxMiddleware(db), handler.InteractionMux)
	router.Get("/whoami", middlewares.TxMiddleware(db), handler.WhoAmI)
	router.Get("/teams", middlewares.TxMiddleware(db), handler.FindTeam)
	router.Get("/teams/:teamId/channels", middlewares.TxMiddleware(db), handler.FindChannels)
	router.Get("/teams/:teamId/channels/:channelId", middlewares.TxMiddleware(db), handler.FindChannelById)
	router.Get("/teams/:teamId/channels/:channelId/messages/latest", middlewares.TxMiddleware(db), handler.FindLatestChannelMessage)
	router.Get("/teams/:teamId/users", middlewares.TxMiddleware(db), handler.FindTeamUsers)
	router.Get("/google/sheets", middlewares.TxMiddleware(db), handler.FindSheet)
}
