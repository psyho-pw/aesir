package google

import (
	"aesir/src/common"
	"aesir/src/common/errors"
	"aesir/src/common/utils"
	"aesir/src/google/dto"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
	"time"
)

//go:generate mockery --name Service --case underscore --inpackage
type Service interface {
	FindSheet() (*sheets.ValueRange, error)
	AppendRow(createVoCDto *dto.CreateVoCDto) error
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
	readRange := "VoC (Raw data)!A2:P"
	resp, err := service.sheetClient.Spreadsheets.Values.Get(service.config.Google.SpreadSheetId, readRange).Do()
	if err != nil {
		logrus.Errorf("Unable to retrieve data from sheet: %v", err)
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return resp, nil
}

func (service *googleService) AppendRow(createVoCDto *dto.CreateVoCDto) error {
	appendRange := "VoC (Raw data)!A:E"

	records := [][]interface{}{{
		time.Now().Format("2006-01-02"),
		createVoCDto.Client.ClientName,
		createVoCDto.IsStakeholder,
		createVoCDto.User.RealName,
		"",
		createVoCDto.VocContent,
	}}
	logrus.Debugf("%v", records)
	row := sheets.ValueRange{MajorDimension: "ROWS", Range: appendRange, Values: records}
	valueInputOption := "RAW"
	insertDataOption := "INSERT_ROWS"

	resp, err := service.sheetClient.Spreadsheets.Values.Append(service.config.Google.SpreadSheetId, appendRange, &row).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Do()
	if err != nil {
		//TODO test error propagation
		utils.PrettyPrint(err)
		return errors.New(fiber.StatusInternalServerError, err.Error())
	}

	utils.PrettyPrint(resp)

	return nil
}

func (service *googleService) WithTx(tx *gorm.DB) Service {

	return service
}
