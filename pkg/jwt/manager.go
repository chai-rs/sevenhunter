package jwt

import (
	"net/http"
	"time"

	errx "github.com/chai-rs/sevenhunter/pkg/error"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManagerConfig struct {
	Secret          string        `split_words:"true" required:"true"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL" split_words:"true" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" split_words:"true" default:"168h"`
}

func (conf *TokenManagerConfig) New() *TokenManager {
	return &TokenManager{config: conf}
}

type TokenManager struct {
	config *TokenManagerConfig
}

func (tm *TokenManager) SignClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tm.config.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (tm *TokenManager) VerifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(tm.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, errx.M(http.StatusUnauthorized, "unauthorized")
}

func (tm *TokenManager) AccessTokenExpiresAt(now time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(now.Add(tm.config.AccessTokenTTL))
}

func (tm *TokenManager) RefreshTokenExpiresAt(now time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(now.Add(tm.config.RefreshTokenTTL))
}
