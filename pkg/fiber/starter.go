package fx

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

const DefaultAddr = ":8080"

var DefaultStartOpts = StartOpts{
	Address:    DefaultAddr,
	ShutdownFn: func() error { return nil },
}

type StartOpts struct {
	Address    string
	ShutdownFn func() error
}

func Start(app *fiber.App, opts ...StartOpts) error {
	app.Use(healthcheck.New())
	app.Use(notfound)

	opt := DefaultStartOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	addr := opt.Address
	if addr == "" {
		addr = DefaultAddr
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		logx.Info().Msg("shutting down...")

		if err := opt.ShutdownFn(); err != nil {
			logx.Error().Err(err).Msg("shutdown error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logx.Error().Err(err).Msg("shutdown error")
		} else {
			logx.Info().Msg("shutdown complete")
		}
	}()

	logx.Info().Str("addr", addr).Msg("starting application...")
	if err := app.Listen(addr); err != nil && err != http.ErrServerClosed {
		return err
	}

	logx.Info().Msg("application stopped")
	return nil
}

func notfound(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(Response{Success: false, Message: "method not found"})
}
