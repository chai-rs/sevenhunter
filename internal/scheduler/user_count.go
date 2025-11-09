package scheduler

import (
	"context"

	"github.com/chai-rs/sevenhunter/internal/model"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
)

type UserCountScheduler struct {
	userRepo model.UserRepo
}

type UserCountSchedulerOpts struct {
	UserRepo model.UserRepo
}

func NewUserCountScheduler(opts UserCountSchedulerOpts) *UserCountScheduler {
	return &UserCountScheduler{
		userRepo: opts.UserRepo,
	}
}

func (s *UserCountScheduler) Run(ctx context.Context) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		logx.Error().Err(err).Msg("failed to get user count")
	}

	logx.Info().Int64("user_count", count).Msg("current user count")
}
