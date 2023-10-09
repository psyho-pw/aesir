package channels

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --case underscore --inpackage
type Repository interface {
	Create(channel Channel) (*Channel, error)
	CreateMany(channels []Channel) ([]Channel, error)
	FindMany() ([]Channel, error)
	FindManyWithMessage() ([]Channel, error)
	FindManyByThreshold(threshold int) ([]Channel, error)
	FindOneBySlackId(slackId string) (*Channel, error)
	FindFirstOne() (*Channel, error)
	UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error)
	UpdateThreshold(threshold int) error
	DeleteOneBySlackId(slackId string) (*Channel, error)
	WithTx(tx *gorm.DB) Repository
}

type channelRepository struct {
	DB *gorm.DB
}

func NewChannelRepository(db *gorm.DB) Repository {
	return &channelRepository{db}
}

var SetRepository = wire.NewSet(NewChannelRepository)

func (repository *channelRepository) Create(channel Channel) (*Channel, error) {
	result := repository.DB.Create(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}

	return &channel, nil
}

func (repository *channelRepository) CreateMany(channels []Channel) ([]Channel, error) {
	result := repository.DB.Create(&channels)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected != int64(len(channels)) {
		return nil, errors.New(fiber.StatusInternalServerError, "affected row count doesn't match")
	}

	return channels, nil
}

func (repository *channelRepository) FindMany() ([]Channel, error) {
	var channels []Channel
	if err := repository.DB.Omit("Message").Order("id desc").Find(&channels).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return nil, errors.New(fiber.StatusServiceUnavailable, "test")

	return channels, nil
}

func (repository *channelRepository) FindManyWithMessage() ([]Channel, error) {
	var channels []Channel
	result := repository.DB.Preload("Message").Find(&channels)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusInternalServerError, result.Error.Error())
	}

	//if err := repository.DB.Select("*").Joins("left join messages on channel.id = messages.channelId").Scan(&channels).Error; err != nil {
	//	return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	//}

	return channels, nil
}

func (repository *channelRepository) FindManyByThreshold(threshold int) ([]Channel, error) {
	var channels []Channel
	result := repository.DB.Joins("Message").Where("timestamp < UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL ? day))", -threshold).Find(&channels)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusInternalServerError, result.Error.Error())
	}

	return channels, nil
}

func (repository *channelRepository) FindOneBySlackId(slackId string) (*Channel, error) {
	var channel Channel
	result := repository.DB.Preload("Message").Where(&Channel{SlackId: slackId}).Find(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &channel, nil
}

func (repository *channelRepository) FindFirstOne() (*Channel, error) {
	var channel Channel
	result := repository.DB.First(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &channel, nil
}

func (repository *channelRepository) UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error) {
	result := repository.DB.Where(&Channel{SlackId: slackId}).Updates(&channel)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}

	return &channel, nil
}

func (repository *channelRepository) UpdateThreshold(threshold int) error {
	result := repository.DB.Where("1 = 1").Updates(Channel{Threshold: threshold})
	if result.Error != nil {
		return errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New(fiber.StatusNotFound, "not affected")
	}

	return nil
}

func (repository *channelRepository) DeleteOneBySlackId(slackId string) (*Channel, error) {
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

func (repository *channelRepository) WithTx(tx *gorm.DB) Repository {
	repository.DB = tx
	return repository
}
