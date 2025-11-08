package config

import "github.com/chai-rs/sevenhunter/pkg/jwt"

type Config struct {
	Auth *jwt.TokenManagerConfig `required:"true"`
}
