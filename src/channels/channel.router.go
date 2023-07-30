package channels

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler Handler) {
	router.Get("", middlewares.TxMiddleware(db), handler.FineMany)
	router.Get("/test", middlewares.TxMiddleware(db), handler.FindManyWithMessage)
	router.Get("/:id", middlewares.TxMiddleware(db), handler.FindOneBySlackId)
}
