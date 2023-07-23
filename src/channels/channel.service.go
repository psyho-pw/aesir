package channels

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

type ChannelService interface {
	Create(channel Channel) (*Channel, error)
	FindMany() ([]Channel, error)
	FindOneBySlackId(slackId string) (*Channel, error)
	UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error)
	DeleteOneBySlackId(slackId string) (*Channel, error)
	WithTx(tx *gorm.DB) ChannelService
}

type channelService struct {
	repository ChannelRepository
}

func NewChannelService(channelRepository ChannelRepository) ChannelService {
	return &channelService{channelRepository}
}

var SetService = wire.NewSet(NewChannelService)

func (service *channelService) Create(channel Channel) (*Channel, error) {
	result, err := service.repository.Create(channel)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *channelService) FindMany() ([]Channel, error) {
	result, err := service.repository.FindMany()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *channelService) FindOneBySlackId(slackId string) (*Channel, error) {
	result, err := service.repository.FindOneBySlackId(slackId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *channelService) UpdateOneBySlackId(slackId string, channel Channel) (*Channel, error) {
	result, err := service.repository.UpdateOneBySlackId(slackId, channel)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *channelService) DeleteOneBySlackId(slackId string) (*Channel, error) {
	result, err := service.repository.DeleteOneBySlackId(slackId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *channelService) WithTx(tx *gorm.DB) ChannelService {
	service.repository = service.repository.WithTx(tx)
	return service
}
