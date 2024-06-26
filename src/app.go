package src

import (
	"aesir/src/channels"
	"aesir/src/clients"
	"aesir/src/common"
	"aesir/src/common/database"
	"aesir/src/common/middlewares"
	"aesir/src/google"
	"aesir/src/messages"
	"aesir/src/slack"
	"aesir/src/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var AppSet = wire.NewSet(
	common.ConfigSet,
	database.DBSet,
	users.SetRepository,
	users.SetService,
	users.SetHandler,
	channels.SetRepository,
	channels.SetService,
	channels.SetHandler,
	messages.SetRepository,
	messages.SetService,
	messages.SetHandler,
	clients.SetRepository,
	clients.SetService,
	clients.SetHandler,
	slack.SetService,
	slack.SetHandler,
	google.SetService,
	NewApp,
)

func NewApp(
	config *common.Config,
	db *gorm.DB,
	userHandler users.Handler,
	channelHandler channels.Handler,
	messageHandler messages.Handler,
	clientHandler clients.Handler,
	slackHandler slack.Handler,
	googleService google.Service,
) *fiber.App {
	app := fiber.New(config.Fiber)

	if !fiber.IsChild() {
		logrus.Debug("Master process init")
	} else {
		logrus.Debug("Child process init")
	}

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Use(helmet.New())
	//app.Use(csrf.New(config.Csrf))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(middlewares.LogMiddleware)

	app.Static("/", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// router setting
	users.NewRouter(v1.Group("/users"), db, userHandler)
	channels.NewRouter(v1.Group("/channels"), db, channelHandler)
	messages.NewRouter(v1.Group("/messages"), db, messageHandler)
	clients.NewRouter(v1.Group("/clients"), db, clientHandler)
	slack.NewRouter(v1.Group("/slack"), db, slackHandler)

	//err := cronService.Start()
	//if err != nil {
	//	panic(err)
	//}

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return app
}
