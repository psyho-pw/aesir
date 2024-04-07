package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func init() {
	appEnv := os.Getenv("APP_ENV")

	currentWorkDirectory, _ := os.Getwd()
	envPath := currentWorkDirectory + `/.env/.env.` + appEnv
	logrus.Infof("envPath: %s", envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env file")
	}

	for _, e := range os.Environ() {
		logrus.Info(e)
	}

}

func main() {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}

	// Set the timezone for the current process
	time.Local = location

	server, _ := New()
	port := os.Getenv("PORT")
	address := func(appEnv string) string {
		if appEnv == "development" {
			return "localhost"
		}
		return ""
	}(os.Getenv("APP_ENV"))

	log.Fatal(server.Listen(fmt.Sprintf("%s:%s", address, port)))
}
