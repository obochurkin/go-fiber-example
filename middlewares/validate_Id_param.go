package middlewares

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/obochurkin/go-fiber-example/errors"
)

func ValidateIdParam() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Next()
		}

		validId, err := doValidation(id)
		if err != nil || validId == 0 {
			return errors.NotFoundError()
		}

		c.Locals("id", validId)

		return c.Next()
	}
} 

func doValidation(id string) (uint, error) {
	trimmedId := strings.TrimSpace(id)
	// check empty
	if trimmedId == "" {
		return 0, errors.NotFoundError()
	}

	// check max length
	if len(trimmedId) > 50 {
		return 0, errors.BadRequestError()
	}

	// check digits only
	matched, err := regexp.MatchString(`^\d+$`, trimmedId)
	if err != nil || !matched {
		return 0, errors.BadRequestError()
	}

	// check valid uint and range
	validId, err := strconv.ParseUint(trimmedId, 10, 32)
	if err != nil || validId == 0 {
		return 0, errors.BadRequestError()
	}

	return uint(validId), nil
}