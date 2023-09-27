package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func PrettyPrint(input interface{}) {
	prettify, _ := json.MarshalIndent(input, "", "  ")
	logrus.Println(string(prettify))
}
