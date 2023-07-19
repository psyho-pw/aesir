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

var GeneralErrorHandler = func(ctx *fiber.Ctx, err error) error {
	logrus.SetOutput(os.Stderr)
	code := fiber.StatusInternalServerError

	var exception *Errors.Error
	if errors.As(err, &exception) {
		code = exception.Code
	}

	tx := ctx.Locals("TX").(*gorm.DB)
	if tx != nil {
		tx.Rollback()
		logrus.Error("Transaction rollback - GeneralErrorHandler")
	}

	logrus.Errorf("%+v", exception)
	logrus.SetOutput(colorable.NewColorableStdout())

	return ctx.Status(code).JSON(exception)
}
