package model

import "github.com/golang-jwt/jwt/v5"

type AuthTokenType string

func (t AuthTokenType) String() string {
	return string(t)
}

const (
	AccessToken  AuthTokenType = "access_token"
	RefreshToken AuthTokenType = "refresh_token"
)

type AccessTokenClaims struct {
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Type  AuthTokenType `json:"type"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	Type AuthTokenType `json:"type"`
	jwt.RegisteredClaims
}

type RegisterOpts struct {
	Name     string
	Email    string
	Password string
}

type LoginOpts struct {
	Email    string
	Password string
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *User
}
