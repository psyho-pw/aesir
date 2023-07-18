package users

import (
	"fiber/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router, db *gorm.DB, handler UserHandler) {
	router.Post("/users", middlewares.TxMiddleware(db), handler.CreateOne)
	router.Get("/users", middlewares.TxMiddleware(db), handler.FindMany)
	router.Get("/users/:id", middlewares.TxMiddleware(db), handler.FindOne)
	router.Patch("/users/:id", middlewares.TxMiddleware(db), handler.UpdateOne)
	router.Delete("/users", middlewares.TxMiddleware(db), handler.DeleteOne)
}
