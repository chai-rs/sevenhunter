package router

import (
	"github.com/chai-rs/sevenhunter/internal/handler"
	"github.com/chai-rs/sevenhunter/internal/repo"
	"github.com/chai-rs/sevenhunter/internal/service"
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type BindAuthOpts struct {
	DB           *mongo.Database
	TokenManager *jwt.TokenManager
}

func BindAuth(group fiber.Router, opts BindAuthOpts) {
	hdl := handler.NewAuthHandler(handler.AuthHandlerOpts{
		Service: service.NewAuthService(&service.AuthServiceOpts{
			TokenManager: opts.TokenManager,
			UserRepo:     repo.NewUserRepo(opts.DB),
		}),
	})

	router := group.Group("/auth")
	router.Post("/login", hdl.Login)
	router.Post("/register", hdl.Register)
	router.Post("/refresh", hdl.RefreshToken)
}
