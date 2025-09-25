package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/obochurkin/go-fiber-example/errors"
)

func GetUsers(c *fiber.Ctx) error {
	return errors.NotFoundError()
	//return c.JSON(fiber.Map{"users": []string{"Alice", "Bob", "Charlie"}})
}