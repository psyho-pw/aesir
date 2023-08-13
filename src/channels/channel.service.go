package channels

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Service --case underscore --inpackage
type Service interface {
	Create(channel Channel) (*Channel, error)
	CreateMany(channels []Channel) ([]Channel, error)
	FindMany() ([]Channel, error)
	FindManyWithMessage() ([]Channel, error)
	FindOneBySlackId(slackId string) (*Channel, error)
	FindFirstOne() (*Channel, error)
	UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error)
	UpdateThreshold(threshold int) error
	DeleteOneBySlackId(slackId string) (*Channel, error)
	WithTx(tx *gorm.DB) Service
}

type channelService struct {
	repository Repository
}

func NewChannelService(channelRepository Repository) Service {
	return &channelService{channelRepository}
}

var SetService = wire.NewSet(NewChannelService)

func (service *channelService) Create(channel Channel) (*Channel, error) {
	return service.repository.Create(channel)
}

func (service *channelService) CreateMany(channels []Channel) ([]Channel, error) {
	return service.repository.CreateMany(channels)
}

func (service *channelService) FindMany() ([]Channel, error) {
	return service.repository.FindMany()
}

func (service *channelService) FindManyWithMessage() ([]Channel, error) {
	return service.repository.FindManyWithMessage()
}

func (service *channelService) FindOneBySlackId(slackId string) (*Channel, error) {
	return service.repository.FindOneBySlackId(slackId)
}

func (service *channelService) FindFirstOne() (*Channel, error) {
	return service.repository.FindFirstOne()
}

func (service *channelService) UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error) {
	return service.repository.UpdateOneBySlackId(slackId, channel)
}

func (service *channelService) UpdateThreshold(threshold int) error {
	return service.repository.UpdateThreshold(threshold)
}

func (service *channelService) DeleteOneBySlackId(slackId string) (*Channel, error) {
	return service.repository.DeleteOneBySlackId(slackId)
}

func (service *channelService) WithTx(tx *gorm.DB) Service {
	service.repository = service.repository.WithTx(tx)
	return service
}
