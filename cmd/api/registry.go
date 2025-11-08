package main

import (
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	MongoDB      *mongo.Client
	TokenManager *jwt.TokenManager
}
