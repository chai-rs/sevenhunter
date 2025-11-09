package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chai-rs/sevenhunter/internal/model"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func Auth(tm *jwt.TokenManager, userRepo model.UserRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			return errx.M(http.StatusBadRequest, "authorization header missing")
		}

		if !strings.HasPrefix(token, "Bearer ") {
			return errx.M(http.StatusBadRequest, "invalid bearer token")
		}

		token = strings.TrimPrefix(token, "Bearer ")
		parsedToken, err := tm.VerifyToken(token)
		if err != nil {
			return err
		}

		bytes, err := json.Marshal(parsedToken.Claims)
		if err != nil {
			return errx.InternalServerError
		}

		var claims model.AccessTokenClaims
		if err := json.Unmarshal(bytes, &claims); err != nil {
			return errx.InternalServerError
		}

		if claims.Type != model.AccessToken {
			return errx.M(http.StatusUnauthorized, "invalid token type")
		}

		userID := claims.Subject
		exist, err := userRepo.ExistsByID(c.Context(), userID)
		if err != nil {
			return err
		}

		if !exist {
			return errx.M(http.StatusUnauthorized, "user not found")
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
