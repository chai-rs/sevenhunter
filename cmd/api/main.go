package main

import (
	"context"

	"github.com/chai-rs/sevenhunter/cmd/api/config"
	_ "github.com/chai-rs/sevenhunter/docs"
	"github.com/chai-rs/sevenhunter/internal/router"
	fx "github.com/chai-rs/sevenhunter/pkg/fiber"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	_ "github.com/chai-rs/sevenhunter/pkg/logger/autoload"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/samber/lo"
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

// @title SevenHunter API
// @version 1.0
// @description API for SevenHunter assignment
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email tk.thanatorn@gmail.com

// @host localhost:8080
// @BasePath /v1/api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler:   fx.ErrorHandler,
		ReadBufferSize: 16384,
	})

	// Common middlewares
	app.Use(cors.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: lo.ToPtr(logx.ConsoleWriter()),
	}))

	// Bind API routes
	bindAPI(app)
	app.Get("/swagger/*", swagger.HandlerDefault)

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
