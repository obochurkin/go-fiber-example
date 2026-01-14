package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/obochurkin/go-fiber-example/handlers"
	"github.com/obochurkin/go-fiber-example/middlewares"
)

func InitRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	//Init Controllers
	userController := &handlers.UserController{}
	authController := &handlers.AuthController{}

	v1.Route("/auth", func(v1 fiber.Router) {

		v1.Post("/login", authController.Login)
	})

	v1.Route("/users", func(v1 fiber.Router) {

		v1.Get("/", userController.GetUsers)
		v1.Get("/:id", middlewares.ValidateIdParam(), userController.GetUserById)
		v1.Post("/register", userController.CreateUser)
	})
}