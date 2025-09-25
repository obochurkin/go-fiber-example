package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/obochurkin/go-fiber-example/handlers"
)

func InitRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Route("/users", func(v1 fiber.Router) {
		v1.Get("/", handlers.GetUsers)
		v1.Post("/", func(c *fiber.Ctx) error {
			type User struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			p := new(User)

			if err := c.BodyParser(p); err != nil {
				return err
			}

			return c.SendStatus(fiber.StatusCreated)
		})
	})
}