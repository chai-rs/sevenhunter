package config

import (
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"github.com/chai-rs/sevenhunter/pkg/mongo"
)

type Config struct {
	Mongo *mongo.Config           `required:"true"`
	Auth  *jwt.TokenManagerConfig `required:"true"`
}
