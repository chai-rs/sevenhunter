package service

import (
	"context"
	"net/http"
	"time"

	"github.com/chai-rs/sevenhunter/internal/model"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	jwtx "github.com/chai-rs/sevenhunter/pkg/jwt"
	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	tokenManager *jwtx.TokenManager
	userRepo     model.UserRepo
}

type AuthServiceOpts struct {
	TokenManager *jwtx.TokenManager
	UserRepo     model.UserRepo
}

func NewAuthService(opts *AuthServiceOpts) *AuthService {
	return &AuthService{
		tokenManager: opts.TokenManager,
		userRepo:     opts.UserRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, opts model.RegisterOpts) (*model.AuthResult, error) {
	exist, _ := s.userRepo.FindByEmail(ctx, opts.Email)
	if exist != nil {
		logx.Error().Msg("user have already exists")
		return nil, errx.M(http.StatusBadRequest, "user have already exists")
	}

	newUser, err := model.NewCreateUser(model.CreateUserOpts(opts))
	if err != nil {
		logx.Error().Err(err).Msg("failed to create user option")
		return nil, err
	}

	user, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		logx.Error().Err(err).Msg("failed to create new user")
		return nil, err
	}

	atk, err := s.generateAccessToken(user)
	if err != nil {
		logx.Error().Err(err).Msg("failed to generate access token")
		return nil, err
	}

	rtk, err := s.generateRefreshToken(user)
	if err != nil {
		logx.Error().Err(err).Msg("failed to generate refresh token")
		return nil, err
	}

	return &model.AuthResult{
		AccessToken:  atk,
		RefreshToken: rtk,
		User:         user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, opts model.LoginOpts) (*model.AuthResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, opts.Email)
	if err != nil {
		logx.Error().Err(err).Msg("failed to find user by email")
		return nil, err
	}

	if err := user.ComparePassword(opts.Password); err != nil {
		logx.Error().Err(err).Msg("invalid user password")
		return nil, err
	}

	atk, err := s.generateAccessToken(user)
	if err != nil {
		logx.Error().Err(err).Msg("failed to generate access token")
		return nil, err
	}

	rtk, err := s.generateRefreshToken(user)
	if err != nil {
		logx.Error().Err(err).Msg("failed to generate refresh token")
		return nil, err
	}

	return &model.AuthResult{
		AccessToken:  atk,
		RefreshToken: rtk,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.AuthResult, error) {
	token, err := s.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.RefreshTokenClaims)
	if !ok {
		logx.Error().Msg("invalid refresh token claims")
		return nil, errx.M(http.StatusUnauthorized, "invalid refresh token claims")
	}

	user, err := s.userRepo.FindByID(ctx, claims.Subject)
	if err != nil {
		logx.Error().Err(err).Msg("failed to find user by id from refresh token")
		return nil, err
	}

	atk, err := s.generateAccessToken(user)
	if err != nil {
		logx.Error().Err(err).Msg("failed to generate access token")
		return nil, err
	}

	return &model.AuthResult{
		AccessToken:  atk,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	return s.tokenManager.SignClaims(&model.AccessTokenClaims{
		Name:  user.Name(),
		Email: user.Email(),
		Type:  model.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   user.ID(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: s.tokenManager.AccessTokenExpiresAt(now),
		},
	})
}

func (s *AuthService) generateRefreshToken(user *model.User) (string, error) {
	now := time.Now()
	return s.tokenManager.SignClaims(&model.RefreshTokenClaims{
		Type: model.RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   user.ID(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: s.tokenManager.RefreshTokenExpiresAt(now),
		},
	})
}
