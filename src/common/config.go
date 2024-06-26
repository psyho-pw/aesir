package common

import (
	"aesir/src/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/wire"
	"github.com/mattn/go-colorable"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type DB struct {
	MariadbHost     string
	MariadbUsername string
	MariadbPassword string
	MariadbDatabase string
	MariadbPort     string
}

type SlackConfig struct {
	AppToken string
	BotToken string
	TeamId   string
}

type DiscordConfig struct {
	WebhookUrl string
}

type GoogleConfig struct {
	CredentialFilePath string
	SpreadSheetId      string
}

type Config struct {
	AppEnv  string
	Port    int
	Fiber   fiber.Config
	DB      DB
	Csrf    csrf.Config
	Slack   SlackConfig
	Discord DiscordConfig
	Google  GoogleConfig
}

func fiberConfig(webhookUrl string) fiber.Config {
	return fiber.Config{
		//Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Fiber v1",
		ErrorHandler:  middlewares.GeneralErrorHandler(webhookUrl),
	}
}

func dbConfig() DB {
	return DB{
		MariadbHost:     os.Getenv("MARIADB_HOST"),
		MariadbUsername: os.Getenv("MARIADB_USERNAME"),
		MariadbPassword: os.Getenv("MARIADB_PASSWORD"),
		MariadbDatabase: os.Getenv("MARIADB_DATABASE"),
		MariadbPort:     os.Getenv("MARIADB_PORT"),
	}
}

type LumberjackHook struct {
	Writer *lumberjack.Logger
}

func (hook *LumberjackHook) Fire(entry *log.Entry) error {
	msg, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(msg))
	return err
}

func (hook *LumberjackHook) Levels() []log.Level {
	return log.AllLevels
}

func NewLumberjackHook(writer *lumberjack.Logger) *LumberjackHook {
	return &LumberjackHook{
		Writer: writer,
	}
}

func loggerConfig() {
	logPath := "./logs/logs.log"
	maxSize := 10
	maxBackups := 90
	maxAge := 1

	logRotation := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,    // 파일 최대 크기 (MB)
		MaxBackups: maxBackups, // 보관할 백업 파일의 최대 개수
		MaxAge:     maxAge,     // 보관 기간 (일)
		Compress:   true,       // 압축 여부
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	log.SetLevel(log.DebugLevel)
	//log.SetReportCaller(true)
	log.SetOutput(colorable.NewColorableStdout())

	logHook := NewLumberjackHook(logRotation)
	log.AddHook(logHook)
}

func csrfConfig() csrf.Config {
	return csrf.Config{
		KeyLookup:      "header:X-Csrf-BotToken", // string in the form of '<source>:<key>' that is used to extract token from the request
		CookieName:     "csrf_",                  // name of the session cookie
		CookieSameSite: "Lax",                    // indicates if CSRF cookie is requested by SameSite
		Expiration:     3 * time.Hour,            // expiration is the duration before CSRF token will expire
		KeyGenerator:   utils.UUID,               // creates a new CSRF token
	}
}

func slackConfig() SlackConfig {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	if appToken == "" {
		log.Error("Missing slack app token")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		log.Error("app token must have the prefix \"xapp-\"")
	}

	if botToken == "" {
		log.Error("Missing slack bot token")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		log.Error("bot token must have the prefix \"xoxb-\"")
	}

	return SlackConfig{
		AppToken: appToken,
		BotToken: botToken,
		TeamId:   os.Getenv("TEAM_ID"),
	}
}

func discordConfig() DiscordConfig {
	return DiscordConfig{
		WebhookUrl: os.Getenv("DISCORD_WEBHOOK_URL"),
	}
}

func googleConfig() GoogleConfig {
	currentWorkDirectory, _ := os.Getwd()

	return GoogleConfig{
		CredentialFilePath: currentWorkDirectory + `/.env/aesir.json`,
		SpreadSheetId:      os.Getenv("SHEET_ID"),
	}
}

func NewConfig() *Config {
	port, parseErr := strconv.Atoi(os.Getenv("PORT"))
	if parseErr != nil {
		panic(parseErr)
	}

	loggerConfig()

	discord := discordConfig()

	var config = Config{
		AppEnv:  os.Getenv("APP_ENV"),
		Port:    port,
		Fiber:   fiberConfig(discord.WebhookUrl),
		DB:      dbConfig(),
		Csrf:    csrfConfig(),
		Slack:   slackConfig(),
		Discord: discordConfig(),
		Google:  googleConfig(),
	}

	return &config
}

var ConfigSet = wire.NewSet(NewConfig)
