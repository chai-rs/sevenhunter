package model

import "context"

type UserRepo interface {
	Count() (int64, error)
	List(ctx context.Context, opts ListUserOpts) ([]User, error)
	Create(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
