package errors

import "github.com/gofiber/fiber/v2"

func InternalError() *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

}