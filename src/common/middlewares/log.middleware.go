package middlewares

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
)

var LogMiddleware = func(c *fiber.Ctx) error {
	logrus.Info(c)
	queryParams := c.Request().URI().QueryArgs().String()
	if strings.EqualFold(queryParams, "") {
		return c.Next()
	}
	logrus.Infof("Query: %s", queryParams)

	var prettyBodyParams interface{}
	if c.Body() == nil {
		return c.Next()
	}
	err := json.Unmarshal(c.Body(), &prettyBodyParams)
	if err != nil {
		logrus.Errorf("Failed to unmarshal body parameters: %v", err)
		return c.Next()
	}

	prettyBody, errIndent := json.MarshalIndent(prettyBodyParams, "", "  ")
	if errIndent != nil {
		logrus.Errorf("Failed to marshal body parameters: %v", err)
		return c.Next()
	}

	logrus.Infof("Body: %s", string(prettyBody))
	return c.Next()
}
