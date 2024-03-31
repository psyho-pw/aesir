package clients

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler Handler) {
	router.Post("", middlewares.TxMiddleware(db), handler.CreateOne)
	router.Get("", middlewares.TxMiddleware(db), handler.FindMany)
	router.Delete("/:id", middlewares.TxMiddleware(db), handler.DeleteOne)
}
