package handler

import "github.com/chai-rs/sevenhunter/internal/model"

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
