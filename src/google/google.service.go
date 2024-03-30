package google

import (
	"aesir/src/common"
	"aesir/src/common/errors"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

type Service interface {
	FindSheet() (*sheets.ValueRange, error)
	WithTx(tx *gorm.DB) Service
}

type googleService struct {
	sheetClient *sheets.Service
	config      *common.Config
}

func NewGoogleService(config *common.Config) Service {
	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(config.Google.CredentialFilePath))
	if err != nil {
		panic(err)
	}

	return &googleService{
		config:      config,
		sheetClient: client,
	}
}

var SetService = wire.NewSet(NewGoogleService)

/*
********** Sheet related services
 */

func (service *googleService) FindSheet() (*sheets.ValueRange, error) {
	readRange := "측정기록 Fix!A2:P"
	resp, err := service.sheetClient.Spreadsheets.Values.Get(service.config.Google.SpreadSheetId, readRange).Do()
	if err != nil {
		logrus.Errorf("Unable to retrieve data from sheet: %v", err)
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}
	return resp, nil
}

func (service *googleService) WithTx(tx *gorm.DB) Service {

	return service
}
