package handler

import (
	"net/http"

	"github.com/chai-rs/sevenhunter/internal/dto"
	"github.com/chai-rs/sevenhunter/internal/model"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	fx "github.com/chai-rs/sevenhunter/pkg/fiber"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service model.UserService
}

type UserHandlerOpts struct {
	Service model.UserService
}

func NewUserHandler(opts UserHandlerOpts) *UserHandler {
	return &UserHandler{
		service: opts.Service,
	}
}

func (h *UserHandler) getUserID(c *fiber.Ctx) (string, error) {
	userID := c.Locals("user_id", "").(string)
	if userID == "" {
		return "", errx.M(http.StatusUnauthorized, "unauthorized")
	}

	return userID, nil
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	var req dto.ListUsersReq
	if err := c.QueryParser(&req); err != nil {
		return err
	}

	users, err := h.service.List(c.Context(), req.Model())
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.NewUsersRespList(users))
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	user, err := h.service.Get(c.Context(), userID)
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.NewUserResp(user))
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	var req dto.UpdateUserReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	user, err := h.service.Update(c.Context(), req.Model(userID))
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.NewUserResp(user))
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Context(), userID); err != nil {
		return err
	}

	return fx.Ok(c)
}

func (h *UserHandler) Count(c *fiber.Ctx) error {
	count, err := h.service.Count(c.Context())
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.CountUsersResp{Count: count})
}
