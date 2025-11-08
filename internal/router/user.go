package router

import (
	"github.com/chai-rs/sevenhunter/internal/handler"
	"github.com/chai-rs/sevenhunter/internal/repo"
	"github.com/chai-rs/sevenhunter/internal/service"
	"github.com/chai-rs/sevenhunter/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type BindUserOpts struct {
	DB           *mongo.Database
	TokenManager *jwt.TokenManager
}

func BindUser(group fiber.Router, opts BindUserOpts) {
	hdl := handler.NewUserHandler(handler.UserHandlerOpts{
		Service: service.NewUserService(service.UserServiceOpts{
			UserRepo: repo.NewUserRepo(opts.DB),
		}),
	})

	router := group.Group("/users")
	router.Get("", hdl.List)
	router.Get("/count", hdl.Count)
	router.Get("/profile", hdl.Get)
	router.Put("/profile", hdl.Update)
	router.Delete("/profile", hdl.Delete)
}
