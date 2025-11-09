package fx

import (
	"net/http"

	errx "github.com/chai-rs/sevenhunter/pkg/error"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// Response represents the standard API response wrapper
type Response struct {
	Success bool   `json:"success" example:"true"`                                       // Indicates if the request was successful
	Message string `json:"message,omitempty" example:"Operation completed successfully"` // Error or informational message
	Result  any    `json:"result,omitempty"`                                             // The actual response data
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
		} else {
			logx.Error().Err(e).Msg("internal server error")
		}

		return c.Status(e.Code).JSON(resp)
	}

	return c.Status(http.StatusInternalServerError).JSON(resp)
}
