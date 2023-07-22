package users

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user User) (*User, error)
	Find() ([]User, error)
	FindOne(id int) (*User, error)
	FindOneBySlackId(id string) (*User, error)
	UpdateOne(id int, user User) (*User, error)
	DeleteOne(id int) (*User, error)
	WithTx(tx *gorm.DB) UserRepository
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

var SetRepository = wire.NewSet(NewUserRepository)

func (repository *userRepository) Create(user User) (*User, error) {
	result := repository.DB.Create(&user)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}
	//if true {
	//return nil, errors.New(fiber.StatusConflict, "transaction error test")
	//}

	return &user, nil
}

func (repository *userRepository) Find() ([]User, error) {
	var users []User
	if err := repository.DB.Order("id desc").Find(&users).Error; err != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, err.Error())
	}

	return users, nil
}

func (repository *userRepository) FindOne(id int) (*User, error) {
	var user User
	result := repository.DB.Find(&user, id)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "Not affected")
	}

	return &user, nil
}

func (repository *userRepository) FindOneBySlackId(id string) (*User, error) {
	var user User
	result := repository.DB.Where(&User{SlackId: id}).Find(&user)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}

	return &user, nil
}

func (repository *userRepository) UpdateOne(id int, user User) (*User, error) {
	result := repository.DB.Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}
	return &user, nil
}

func (repository *userRepository) DeleteOne(id int) (*User, error) {
	var user User
	result := repository.DB.Delete(&user, id)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(fiber.StatusNotFound, "not affected")
	}
	return &user, nil
}

func (repository *userRepository) WithTx(tx *gorm.DB) UserRepository {
	repository.DB = tx
	return repository
}
