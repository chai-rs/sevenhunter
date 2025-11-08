package fx

import (
	"github.com/gofiber/fiber/v2"
)

const DefaultAddr = ":8080"

func Start(app *fiber.App, addrs ...string) {
	addr := DefaultAddr
	if len(addrs) > 0 {
		addr = addrs[0]
	}

	app.Listen(addr)
}
