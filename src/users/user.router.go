package users

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler UserHandler) {
	router.Get("", middlewares.TxMiddleware(db), handler.FindMany)
	router.Get("/:id", middlewares.TxMiddleware(db), handler.FindOne)
}
