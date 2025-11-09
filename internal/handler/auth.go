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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterReq true "Registration details"
// @Success 201 {object} fx.Response{result=dto.AuthResp} "User successfully registered"
// @Failure 400 {object} fx.Response "Invalid request body or validation error"
// @Failure 409 {object} fx.Response "User already exists"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginReq true "Login credentials"
// @Success 200 {object} fx.Response{result=dto.AuthResp} "Successfully authenticated"
// @Failure 400 {object} fx.Response "Invalid request body or validation error"
// @Failure 401 {object} fx.Response "Invalid credentials"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /auth/login [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenReq true "Refresh token"
// @Success 200 {object} fx.Response{result=dto.AuthResp} "New tokens generated successfully"
// @Failure 400 {object} fx.Response "Invalid request body or validation error"
// @Failure 401 {object} fx.Response "Invalid or expired refresh token"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /auth/refresh [post]
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
