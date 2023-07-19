package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TxMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()
		logrus.Debug("%+v", tx)
		logrus.Fatalf("Transaction start")
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("%+v", r)
				tx.Rollback()
				logrus.Error("Transaction rollback")
				_ = c.Status(fiber.StatusInternalServerError).SendString("internal server error")
			}
			logrus.Debug("Transaction end")
		}()
		c.Locals("TX", tx)
		return c.Next()
	}
}
