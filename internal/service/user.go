package service

import (
	"context"

	"github.com/chai-rs/sevenhunter/internal/model"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
)

type UserService struct {
	userRepo model.UserRepo
}

type UserServiceOpts struct {
	UserRepo model.UserRepo
}

func NewUserService(opts UserServiceOpts) *UserService {
	return &UserService{
		userRepo: opts.UserRepo,
	}
}

var _ model.UserService = (*UserService)(nil)

func (s *UserService) Count(ctx context.Context) (int64, error) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		logx.Error().Err(err).Msg("failed to count the users")
		return 0, err
	}

	return count, nil
}

func (s *UserService) List(ctx context.Context, opts model.ListUserOpts) ([]model.User, error) {
	users, err := s.userRepo.List(ctx, opts)
	if err != nil {
		logx.Error().Err(err).Msg("failed to list the users")
		return nil, err
	}

	return users, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		logx.Error().Err(err).Msgf("failed to get the user with id: %s", id)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, opts model.UpdateUserOpts) (*model.User, error) {
	user, err := s.Get(ctx, opts.ID)
	if err != nil {
		logx.Error().Err(err).Msgf("failed to get the user with id: %s for update", opts.ID)
		return nil, err
	}

	if err := user.Update(opts); err != nil {
		logx.Error().Err(err).Msgf("failed to update the user with id: %s", opts.ID)
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		logx.Error().Err(err).Msgf("failed to save the updated user with id: %s", opts.ID)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		logx.Error().Err(err).Msgf("failed to delete the user with id: %s", id)
		return err
	}

	return nil
}
