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
)

const DefaultAddr = ":8080"

func Start(app *fiber.App, addrs ...string) error {
	addr := DefaultAddr
	if len(addrs) > 0 {
		addr = addrs[0]
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		logx.Info().Msg("shutting down...")

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
