package config

import (
	"github.com/chai-rs/sevenhunter/pkg/config"
	_ "github.com/joho/godotenv/autoload"
)

func Read() *Config {
	return config.MustNew[Config]("")
}
