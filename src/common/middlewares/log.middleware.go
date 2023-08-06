package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/url"
)

func printQuery(param map[string]string) {
	if len(param) == 0 {
		return
	}

	queryJSON, err := json.MarshalIndent(param, "", "  ")
	if err != nil {
		logrus.Errorf("Failed to marshal query parameters: %v", err)
		return
	}

	logrus.Infof("Query: %s", string(queryJSON))
}

func printBody(params []byte) {
	if params == nil {
		return
	}

	var prettyBodyParams interface{}

	logrus.Infof("%+v", fmt.Sprintf("%s", params))
	err := json.Unmarshal(params, &prettyBodyParams)
	if err != nil {
		values, parseErr := url.ParseQuery(string(params))
		if parseErr != nil {
			logrus.Errorf("Failed to parse body parameters: %v", parseErr)
			return
		}
		obj := map[string]string{}
		for k, v := range values {
			if len(v) > 0 {
				obj[k] = v[0]
			}
		}

		prettyBody, marshalErr := json.MarshalIndent(obj, "", "  ")
		if marshalErr != nil {
			logrus.Errorf("Failed to marshal body parameters: %v", err)
			return
		}

		logrus.Infof("Body: %s", string(prettyBody))
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
	printQuery(c.Queries())
	printBody(c.Body())

	return c.Next()
}
