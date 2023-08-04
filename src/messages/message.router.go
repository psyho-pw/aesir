package messages

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler Handler) {
	router.Get("", middlewares.TxMiddleware(db), handler.FindMany)
}
