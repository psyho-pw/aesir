package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func init() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "development" {
		return
	}

	currentWorkDirectory, _ := os.Getwd()
	err := godotenv.Load(currentWorkDirectory + `/.env/.env.` + appEnv)
	if err != nil {
		panic("Error loading .env file")
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
