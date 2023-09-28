package utils

import (
	"github.com/sirupsen/logrus"
	"time"
)

func Timer() func() {
	name := CallerName(1)
	start := time.Now()
	return func() {
		logrus.Infof("%s took %v\n", name, time.Since(start))
	}
}
