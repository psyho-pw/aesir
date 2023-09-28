package errors

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func Report(webhookUrl string, exception *Error) error {
	credentialsResponse, credentialsErr := http.Get(webhookUrl)
	defer func(response *http.Response) {
		if r := recover(); r != nil {
			logrus.Errorf("Reporter credential error")
		}
	}(credentialsResponse)

	if credentialsErr != nil {
		panic(credentialsErr)
	}

	data, readErr := io.ReadAll(credentialsResponse.Body)
	if readErr != nil {
		panic(readErr)
	}

	logrus.Printf("%+v", data)
	return nil
}
