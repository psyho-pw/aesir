package main

import (
	"aesir/src/common/utils"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	appEnv := os.Getenv("APP_ENV")

	currentWorkDirectory, _ := os.Getwd()
	envPath := currentWorkDirectory + `/.env/.env.` + appEnv
	err := godotenv.Load(envPath)
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

	env := map[string]string{}
	for _, e := range os.Environ() {
		split := strings.Split(e, "=")
		env[split[0]] = split[1]
	}
	utils.PrettyPrint(env)
	log.Fatal(server.Listen(fmt.Sprintf("%s:%s", address, port)))
}
