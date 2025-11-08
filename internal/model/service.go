package model

import "context"

type AuthService interface {
	Register(ctx context.Context, opts RegisterOpts) (*AuthResult, error)
	Login(ctx context.Context, opts LoginOpts) (*AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
}

type UserService interface {
	List(ctx context.Context, opts ListUserOpts) ([]User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, opts UpdateUserOpts) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}
