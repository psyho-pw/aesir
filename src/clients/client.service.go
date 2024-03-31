package clients

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Service --case underscore --inpackage
type Service interface {
	CreateOne(client *Client) (*Client, error)
	FindMany() ([]Client, error)
	DeleteOne(id int) (*Client, error)
	WithTx(tx *gorm.DB) Service
}

type clientService struct {
	repository Repository
}

func NewClientService(clientRepository Repository) Service {
	return &clientService{repository: clientRepository}
}

var SetService = wire.NewSet(NewClientService)

func (service *clientService) CreateOne(client *Client) (*Client, error) {
	return service.repository.Create(*client)
}

func (service *clientService) FindMany() ([]Client, error) {
	return service.repository.FindMany()
}

func (service *clientService) DeleteOne(id int) (*Client, error) {
	return service.repository.DeleteOne(id)
}

func (service *clientService) WithTx(tx *gorm.DB) Service {
	service.repository = service.repository.WithTx(tx)
	return service
}
