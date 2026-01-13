package handlers

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/obochurkin/go-fiber-example/dtos"
	"github.com/obochurkin/go-fiber-example/errors"
	"github.com/obochurkin/go-fiber-example/repositories"
)

type UserController struct{}

var usersRepository = &repositories.UsersRepository{}

func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	users, err := usersRepository.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (uc *UserController) GetUserById(c *fiber.Ctx) error {
	userID := c.Locals("id").(uint)
	user, err := usersRepository.FindByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	userInput := dtos.CreateUserDTO{}
	if err := c.BodyParser(&userInput); err != nil {
		return errors.BadRequestError()
	}

	count ,err := usersRepository.FindByEmail(userInput.Email)
	if err != nil {
		log.Errorf("Error checking existing email: %v", err)
		return errors.InternalError()
	}

	if count > 0 {
		return c.SendStatus(fiber.StatusCreated) // to avoid email enumeration
	}

	salt, err := GenerateSalt(userInput.Email)
	if err != nil {
		log.Errorf("Error generating salt: %v", err)
		return errors.InternalError()
	}	

	stringifiedSalt := base64.StdEncoding.EncodeToString(salt)
	
	hashedPassword, err := HashPassword(userInput.Password, stringifiedSalt)
	if err !=nil {
		return errors.InternalError()
	}

	userInput.Password = hashedPassword

	if err := usersRepository.Create(userInput, stringifiedSalt); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

// func (uc *UserController) UpdateUser(c *fiber.Ctx) error {

// }

