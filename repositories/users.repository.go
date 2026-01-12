package repositories

import (
	"github.com/obochurkin/go-fiber-example/database"
	"github.com/obochurkin/go-fiber-example/dtos"
	"github.com/obochurkin/go-fiber-example/errors"
	"github.com/obochurkin/go-fiber-example/models"
)

type UsersRepository struct {}

func(ur *UsersRepository) FindAll() ([]models.User, error) {
		users := []models.User{}
	if err := database.Instance.DB.Find(&users).Error; err != nil {
		return nil, errors.InternalError()
	}
	return users, nil
}

func(ur *UsersRepository) FindByID(id uint) (models.User, error) {
	user := models.User{}
	if err := database.Instance.DB.First(&user, id).Error; err != nil {
		return models.User{}, errors.NotFoundError()
	}
	return user, nil
}

func(ur *UsersRepository) Create(user dtos.CreateUserDTO, salt string) error {
	newUser := models.User{
		Email:    user.Email,
		Password: user.Password,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Salt:     salt,
	}

	if err := database.Instance.DB.Create(&newUser).Error; err != nil {
		return errors.InternalError()
	}
	return nil
}

func(ur *UsersRepository) FindByEmail(email string) (int64, error) {
	count := int64(0)
	if err := database.Instance.DB.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func(ur *UsersRepository) Update(user interface{}) {}

func(ur *UsersRepository) Delete(id uint) {}