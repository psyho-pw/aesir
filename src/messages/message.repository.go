package messages

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --case underscore --inpackage
type Repository interface {
	Create(message Message) (*Message, error)
	FindMany() ([]Message, error)
	WithTx(tx *gorm.DB) Repository
}

type messageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) Repository {
	return &messageRepository{db}
}

var SetRepository = wire.NewSet(NewMessageRepository)

func (repository *messageRepository) FindMany() ([]Message, error) {
	var messages []Message
	if err := repository.DB.Order("id desc").Find(&messages).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return messages, nil
}

func (repository *messageRepository) Create(message Message) (*Message, error) {
	result := repository.DB.Create(&message)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &message, nil
}

func (repository *messageRepository) WithTx(tx *gorm.DB) Repository {
	repository.DB = tx
	return repository
}
