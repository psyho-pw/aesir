package middlewares

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
)

func printQuery(param string) {
	if strings.EqualFold(param, "") {
		return
	}
	logrus.Infof("Query: %s", param)
}

func printBody(params []byte) {
	if params == nil {
		return
	}

	var prettyBodyParams interface{}
	err := json.Unmarshal(params, &prettyBodyParams)
	if err != nil {
		logrus.Errorf("Failed to unmarshal body parameters: %v", err)
		return
	}

	prettyBody, errIndent := json.MarshalIndent(prettyBodyParams, "", "  ")
	if errIndent != nil {
		logrus.Errorf("Failed to marshal body parameters: %v", err)
		return
	}

	logrus.Infof("Body: %s", string(prettyBody))
}

var LogMiddleware = func(c *fiber.Ctx) error {
	logrus.Info(c)
	printQuery(c.Request().URI().QueryArgs().String())
	printBody(c.Body())

	return c.Next()
}
