package errors

import "github.com/gofiber/fiber/v2"

func ForbiddenError() *fiber.Error {
	return fiber.NewError(fiber.StatusForbidden, "Forbidden")

}