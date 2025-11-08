package main

import (
	"context"

	"github.com/chai-rs/sevenhunter/cmd/api/config"
	"github.com/chai-rs/sevenhunter/internal/router"
	fx "github.com/chai-rs/sevenhunter/pkg/fiber"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	_ "github.com/chai-rs/sevenhunter/pkg/logger/autoload"
	"github.com/gofiber/fiber/v2"
)

var (
	registry *Registry
	conf     *config.Config
)

func init() {
	ctx := context.Background()
	conf = config.Read()
	registry = &Registry{
		MongoDB:      conf.Mongo.MustNew(ctx),
		TokenManager: conf.Auth.New(),
	}
}

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: fx.ErrorHandler,
	})

	// Bind API routes
	bindAPI(app)

	// Start the server
	if err := fx.Start(app); err != nil {
		logx.Error().Err(err).Msg("failed to start application")
		return
	}
}

func bindAPI(app *fiber.App) {
	api := app.Group("/v1/api")
	db := registry.MongoDB.Database("sevenhunter")

	// Auth
	router.BindAuth(api, router.BindAuthOpts{
		DB:           db,
		TokenManager: registry.TokenManager,
	})

	// User
	router.BindUser(api, router.BindUserOpts{
		DB:           db,
		TokenManager: registry.TokenManager,
	})
}
