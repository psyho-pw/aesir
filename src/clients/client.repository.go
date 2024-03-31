package clients

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --case underscore --inpackage
type Repository interface {
	Create(client Client) (*Client, error)
	FindMany() ([]Client, error)
	DeleteOne(id int) (*Client, error)
	WithTx(tx *gorm.DB) Repository
}

type clientRepository struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) Repository {
	return &clientRepository{db}
}

var SetRepository = wire.NewSet(NewClientRepository)

func (repository *clientRepository) Create(client Client) (*Client, error) {
	result := repository.DB.Create(&client)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &client, nil
}

func (repository *clientRepository) FindMany() ([]Client, error) {
	var clients []Client
	if err := repository.DB.Order("id desc").Find(&clients).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return clients, nil
}

func (repository *clientRepository) DeleteOne(id int) (*Client, error) {
	var client Client
	result := repository.DB.Delete(&client, id)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &client, nil
}

func (repository *clientRepository) WithTx(tx *gorm.DB) Repository {
	repository.DB = tx
	return repository
}
