package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/compress"

	"github.com/obochurkin/go-fiber-example/router"
)

const idleTimeout = 5 * time.Second

func main() {
	app := Init()
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	app.Use(requestid.New())
	app.Use(etag.New())
	app.Use(compress.New())

	// TODO: add DB connection

	// Init Routes
	router.InitRoutes(app)



	// Graceful shutdown
	go func() {
		if err := app.Listen(":4003"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nGracefully shutting down...")
	_ = app.Shutdown()
	fmt.Println("Cleanup resources...")

	// dbClose() etc goes here

	fmt.Println("Fiber was successfuly shutdown.")
}

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,

		//
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Return JSON error response
			return ctx.Status(code).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})

	return app
}
