package handler

import (
	"github.com/chai-rs/sevenhunter/internal/dto"
	"github.com/chai-rs/sevenhunter/internal/model"
	fx "github.com/chai-rs/sevenhunter/pkg/fiber"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service model.AuthService
}

type AuthHandlerOpts struct {
	Service model.AuthService
}

func NewAuthHandler(opts AuthHandlerOpts) *AuthHandler {
	return &AuthHandler{
		service: opts.Service,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	result, err := h.service.Register(c.Context(), *req.Model())
	if err != nil {
		return err
	}

	return fx.Created(c, dto.NewAuthResp(result))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	result, err := h.service.Login(c.Context(), *req.Model())
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.NewAuthResp(result))
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	result, err := h.service.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.NewAuthResp(result))
}
