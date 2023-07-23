package users

import (
	"aesir/src/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user User) (*User, error)
	CreateMany(users []User) ([]User, error)
	Find() ([]User, error)
	FindOne(id int) (*User, error)
	FindOneBySlackId(id string) (*User, error)
	UpdateOne(id int, user User) (*User, error)
	DeleteOne(id int) (*User, error)
	WithTx(tx *gorm.DB) Repository
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db}
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

func (repository *userRepository) CreateMany(users []User) ([]User, error) {
	result := repository.DB.Create(&users)
	if result.Error != nil {
		return nil, errors.New(fiber.StatusServiceUnavailable, result.Error.Error())
	}
	if result.RowsAffected != int64(len(users)) {
		return nil, errors.New(fiber.StatusInternalServerError, "affected row count doesn't match")
	}

	return users, nil
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
	if result.RowsAffected == 0 {
		return nil, nil
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

func (repository *userRepository) WithTx(tx *gorm.DB) Repository {
	repository.DB = tx
	return repository
}
