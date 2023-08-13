package users

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate mockery --name Service --case underscore --inpackage
type Service interface {
	CreateOne(*User) (*User, error)
	CreateMany([]User) ([]User, error)
	FindMany() ([]User, error)
	FindOne(id int) (*User, error)
	FindOneBySlackId(id string) (*User, error)
	UpdateOne(id int, user *User) (*User, error)
	UpdateIsManager(ids []int) error
	DeleteOne(id int) (*User, error)
	WithTx(tx *gorm.DB) Service
}

type userService struct {
	repository Repository
}

func NewUserService(userRepository Repository) Service {
	return &userService{repository: userRepository}
}

var SetService = wire.NewSet(NewUserService)

func (service *userService) CreateOne(user *User) (*User, error) {
	return service.repository.Create(*user)

}

func (service *userService) CreateMany(users []User) ([]User, error) {
	return service.repository.CreateMany(users)
}

func (service *userService) FindMany() ([]User, error) {
	return service.repository.Find()
}

func (service *userService) FindOne(id int) (*User, error) {
	return service.repository.FindOne(id)
}

func (service *userService) FindOneBySlackId(id string) (*User, error) {
	return service.repository.FindOneBySlackId(id)
}

func (service *userService) UpdateOne(id int, user *User) (*User, error) {
	return service.repository.UpdateOne(id, *user)
}

func (service *userService) UpdateIsManager(ids []int) error {
	setMangersErr := service.repository.SetManagersByUserIds(ids)
	if setMangersErr != nil {
		return setMangersErr
	}
	removeManagersErr := service.repository.RemoveManagersByUserIds(ids)
	if removeManagersErr != nil {
		return removeManagersErr
	}

	return nil
}

func (service *userService) DeleteOne(id int) (*User, error) {
	return service.repository.DeleteOne(id)
}

func (service *userService) WithTx(tx *gorm.DB) Service {
	service.repository = service.repository.WithTx(tx)
	return service
}
