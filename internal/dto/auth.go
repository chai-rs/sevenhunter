package dto

import "github.com/chai-rs/sevenhunter/internal/model"

type LoginReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=200"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (r *LoginReq) Model() *model.LoginOpts {
	return &model.LoginOpts{
		Email:    r.Email,
		Password: r.Password,
	}
}

type RegisterReq struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email,min=5,max=200"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (r *RegisterReq) Model() *model.RegisterOpts {
	return &model.RegisterOpts{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserResp `json:"user"`
}

func NewAuthResp(m *model.AuthResult) *AuthResp {
	if m == nil {
		return nil
	}

	return &AuthResp{
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		User:         NewUserResp(m.User),
	}
}
