package fx

import (
	"net/http"

	errx "github.com/chai-rs/sevenhunter/pkg/error"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func Ok(c *fiber.Ctx, result ...any) error {
	resp := Response{
		Success: true,
	}

	if len(result) > 0 {
		resp.Result = result[0]
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func Created(c *fiber.Ctx, result ...any) error {
	resp := Response{
		Success: true,
	}

	if len(result) > 0 {
		resp.Result = result[0]
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	resp := Response{
		Success: false,
		Message: errx.InternalServerError.Error(),
	}

	switch e := err.(type) {
	case *errx.Error:
		if e.Code != http.StatusInternalServerError {
			resp.Message = e.Message
		}

		return c.Status(e.Code).JSON(resp)
	}

	return c.Status(http.StatusInternalServerError).JSON(resp)
}
