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
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return "", errx.M(http.StatusUnauthorized, "unauthorized")
	}

	return userID, nil
}

// List godoc
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cursor query string false "Pagination cursor"
// @Param limit query int false "Number of items per page (max 100)" default(10)
// @Param sort_asc query bool false "Sort in ascending order" default(false)
// @Success 200 {object} fx.Response{result=[]dto.UserResp} "Successfully retrieved users list"
// @Failure 401 {object} fx.Response "Unauthorized - invalid or missing token"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /users [get]
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

// Get godoc
// @Summary Get current user profile
// @Description Retrieve the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fx.Response{result=dto.UserResp} "Successfully retrieved user profile"
// @Failure 401 {object} fx.Response "Unauthorized - invalid or missing token"
// @Failure 404 {object} fx.Response "User not found"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /users/profile [get]
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

// Update godoc
// @Summary Update current user profile
// @Description Update the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserReq true "Updated user information"
// @Success 200 {object} fx.Response{result=dto.UserResp} "Successfully updated user profile"
// @Failure 400 {object} fx.Response "Invalid request body or validation error"
// @Failure 401 {object} fx.Response "Unauthorized - invalid or missing token"
// @Failure 404 {object} fx.Response "User not found"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /users/profile [put]
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

// Delete godoc
// @Summary Delete current user account
// @Description Delete the authenticated user's account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fx.Response "Successfully deleted user account"
// @Failure 401 {object} fx.Response "Unauthorized - invalid or missing token"
// @Failure 404 {object} fx.Response "User not found"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /users/profile [delete]
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

// Count godoc
// @Summary Get total user count
// @Description Retrieve the total number of registered users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fx.Response{result=dto.CountUsersResp} "Successfully retrieved user count"
// @Failure 401 {object} fx.Response "Unauthorized - invalid or missing token"
// @Failure 500 {object} fx.Response "Internal server error"
// @Router /users/count [get]
func (h *UserHandler) Count(c *fiber.Ctx) error {
	count, err := h.service.Count(c.Context())
	if err != nil {
		return err
	}

	return fx.Ok(c, dto.CountUsersResp{Count: count})
}
