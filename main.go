package main

import (
	"aesir/src"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	server, _ := src.New()
	port := os.Getenv("PORT")
	address := func(appEnv string) string {
		if appEnv == "development" {
			return "localhost"
		}
		return ""
	}(os.Getenv("APP_ENV"))

	log.Fatal(server.Listen(fmt.Sprintf("%s:%s", address, port)))
}
