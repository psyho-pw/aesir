package database

import (
	"aesir/src/channels"
	"aesir/src/clients"
	"aesir/src/common"
	"aesir/src/messages"
	"aesir/src/users"
	"fmt"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(config *common.Config) *gorm.DB {
	var datetimePrecision = 2

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB.MariadbUsername,
		config.DB.MariadbPassword,
		config.DB.MariadbHost,
		config.DB.MariadbPort,
		config.DB.MariadbDatabase,
	)

	connection, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                      dsn,
			DefaultStringSize:        256,
			DefaultDatetimePrecision: &datetimePrecision,
		}),
		&gorm.Config{
			FullSaveAssociations: true,
		},
	)

	if err != nil {
		panic(err)
	}

	migrationError := connection.AutoMigrate(
		users.User{},
		messages.Message{},
		channels.Channel{},
		clients.Client{},
	)
	if migrationError != nil {
		panic(err)
	}

	if config.AppEnv == "development" {
		connection = connection.Debug()
	}

	return connection
}

var DBSet = wire.NewSet(NewDB)
