package errors

import "github.com/gofiber/fiber/v2"

func BadRequestError() *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, "Bad Request")

}