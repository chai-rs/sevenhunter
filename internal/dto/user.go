package dto

import (
	"github.com/chai-rs/sevenhunter/internal/model"
	"github.com/samber/lo"
)

type UserResp struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

func NewUserResp(m *model.User) *UserResp {
	if m == nil {
		return nil
	}

	return &UserResp{
		ID:        m.ID(),
		Name:      m.Email(),
		Email:     m.Email(),
		CreatedAt: m.CreatedAt().UnixMilli(),
	}
}

func NewUsersRespList(users []model.User) []UserResp {
	return lo.Map(users, func(item model.User, index int) UserResp {
		return *NewUserResp(&item)
	})
}

type ListUsersReq struct {
	Cursor  string `query:"cursor"`
	Limit   int    `query:"limit"`
	SortAsc bool   `query:"sort_asc"`
}

func (r *ListUsersReq) Model() model.ListUserOpts {
	opts := model.ListUserOpts{}
	opts.Cursor = r.Cursor
	opts.SortAsc = r.SortAsc

	if r.Limit <= 0 {
		opts.Limit = 10
	} else if r.Limit > 100 {
		opts.Limit = 100
	}

	return opts
}

type UpdateUserReq struct {
	Name  string `json:"name" validate:"required,min=2,max=32"`
	Email string `json:"email" validate:"required,email,min=5,max=200"`
}

func (r *UpdateUserReq) Model(id string) model.UpdateUserOpts {
	return model.UpdateUserOpts{
		ID:    id,
		Name:  r.Name,
		Email: r.Email,
	}
}

type CountUsersResp struct {
	Count int64 `json:"count"`
}
