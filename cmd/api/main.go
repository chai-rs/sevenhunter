package main

import (
	fx "github.com/chai-rs/sevenhunter/pkg/fiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: fx.ErrorHandler,
	})

	fx.Start(app)
}
