package middlewares

import (
	Errors "aesir/src/common/errors"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

func GeneralErrorHandler(webhookUrl string) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		logrus.SetOutput(os.Stderr)
		code := fiber.StatusInternalServerError

		var exception *Errors.Error
		if errors.As(err, &exception) {
			code = exception.Code
		}

		tx := ctx.Locals("TX")
		if tx != nil {
			tx.(*gorm.DB).Rollback()
			logrus.Error("Transaction rollback - GeneralErrorHandler")
		}

		logrus.Errorf("exception: %+v", exception)

		// error reporting
		reportErr := Errors.Report(webhookUrl, exception)
		if reportErr != nil {
			logrus.Errorf("report service unavailable, %+v", reportErr)
		}

		logrus.SetOutput(colorable.NewColorableStdout())

		return ctx.Status(code).JSON(exception)
	}
}
