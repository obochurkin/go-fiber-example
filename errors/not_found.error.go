package errors

import "github.com/gofiber/fiber/v2"

func NotFoundError() *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound, "Resource not found")

}