package messages

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Service --case underscore --inpackage
type Service interface {
	FindMany() ([]Message, error)
	WithTx(tx *gorm.DB) Service
}

type messageService struct {
	repository Repository
}

func NewMessageService(messageRepository Repository) Service {
	return &messageService{messageRepository}
}

var SetService = wire.NewSet(NewMessageService)

func (service *messageService) FindMany() ([]Message, error) {
	return service.repository.FindMany()
}

func (service *messageService) WithTx(tx *gorm.DB) Service {
	service.repository = service.repository.WithTx(tx)
	return service
}
