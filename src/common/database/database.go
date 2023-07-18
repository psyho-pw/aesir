package database

import (
	"fiber/src/common"
	"fiber/src/users"
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

	connection, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        256,
		DefaultDatetimePrecision: &datetimePrecision,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = connection.AutoMigrate(users.User{})
	if err != nil {
		panic(err)
	}

	if config.AppEnv == "development" {
		connection = connection.Debug()
	}

	return connection
}

var DBSet = wire.NewSet(NewDB)
