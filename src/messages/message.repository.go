package messages

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --case underscore --inpackage
type Repository interface {
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

func (repository *messageRepository) WithTx(tx *gorm.DB) Repository {
	repository.DB = tx
	return repository
}
