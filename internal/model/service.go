package model

import "context"

type AuthService interface {
	Register(ctx context.Context)
	Login(ctx context.Context)
	Refresh(ctx context.Context)
}

type UserService interface {
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}
