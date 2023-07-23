package channels

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type ChannelRepository interface {
	Create(channel Channel) (*Channel, error)
	FindMany() ([]Channel, error)
	FindOneBySlackId(slackId string) (*Channel, error)
	UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error)
	DeleteOneBySlackId(slackId string) (*Channel, error)
	WithTx(tx *gorm.DB) ChannelRepository
}

type channelRepository struct {
	DB *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db}
}

var SetRepository = wire.NewSet(NewChannelRepository)

func (repository channelRepository) Create(channel Channel) (*Channel, error) {
	result := repository.DB.Create(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &channel, nil
}

func (repository channelRepository) FindMany() ([]Channel, error) {
	var channels []Channel
	if err := repository.DB.Order("id desc").Find(&channels).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return channels, nil
}

func (repository channelRepository) FindOneBySlackId(slackId string) (*Channel, error) {
	var channel Channel
	if err := repository.DB.Where(&Channel{SlackId: slackId}).Find(&channel).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return &channel, nil
}

func (repository channelRepository) UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error) {
	result := repository.DB.Where(&Channel{SlackId: slackId}).Updates(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &channel, nil
}

func (repository channelRepository) DeleteOneBySlackId(slackId string) (*Channel, error) {
	var channel Channel
	result := repository.DB.Where(&Channel{SlackId: slackId}).Delete(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &channel, nil
}

func (repository channelRepository) WithTx(tx *gorm.DB) ChannelRepository {
	repository.DB = tx
	return repository
}
