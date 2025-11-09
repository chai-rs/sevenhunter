package config

import (
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"github.com/chai-rs/sevenhunter/pkg/mongo"
)

type Config struct {
	Mongo *mongo.Config           `required:"true"`
	Auth  *jwt.TokenManagerConfig `required:"true"`
	App   *AppConfig              `required:"true"`
}

type AppConfig struct {
	Port               string `env:"APP_PORT" default:"8080"`
	CorsAllowedOrigins string `env:"APP_CORS_ALLOWED_ORIGINS" default:"*"`
	CorsAllowedMethods string `env:"APP_CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
}

func (c AppConfig) Address() string {
	return ":" + c.Port
}
