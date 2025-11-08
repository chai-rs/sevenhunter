package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chai-rs/sevenhunter/internal/model"
	"github.com/chai-rs/sevenhunter/internal/model/mocks"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	jwtx "github.com/chai-rs/sevenhunter/pkg/jwt"
	. "github.com/chai-rs/sevenhunter/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// createTestTokenManager creates a TokenManager instance for testing
func createTestTokenManager() *jwtx.TokenManager {
	config := &jwtx.TokenManagerConfig{
		Secret:          "test-secret-key-for-testing-purposes",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 168 * time.Hour,
	}
	return config.New()
}

// createTestUserWithPassword creates a test user with a known password for testing
func createTestUserWithPassword(id, name, email, password string) *model.User {
	user, err := model.NewCreateUser(model.CreateUserOpts{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		panic("Failed to create test user: " + err.Error())
	}

	userWithID, err := model.NewUser(model.UserOpts{
		ID:             id,
		Name:           user.Name(),
		Email:          user.Email(),
		HashedPassword: user.HashedPassword(),
		CreatedAt:      user.CreatedAt(),
	})
	if err != nil {
		panic("Failed to create test user with ID: " + err.Error())
	}

	return userWithID
}

func TestAuthService_Register(t *testing.T) {
	type Testcase struct {
		name     string
		input    model.RegisterOpts
		arrange  ArrangeFn[*AuthService, model.RegisterOpts]
		validate func(t *testing.T, result *model.AuthResult, err error)
		isError  bool
	}

	testcases := []Testcase{
		{
			name: "register user successfully",
			input: model.RegisterOpts{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.RegisterOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				// User doesn't exist
				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, mongo.ErrNoDocuments)

				// Create user successfully
				repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Name() == input.Name && u.Email() == input.Email
				})).RunAndReturn(func(ctx context.Context, u *model.User) (*model.User, error) {
					// Return user with ID set
					return createTestUserWithPassword("123", u.Name(), u.Email(), input.Password), nil
				})
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEmpty(t, result.AccessToken)
				require.NotEmpty(t, result.RefreshToken)
				require.NotNil(t, result.User)
				require.Equal(t, "John Doe", result.User.Name())
				require.Equal(t, "john@example.com", result.User.Email())
			},
			isError: false,
		},
		{
			name: "register fails when user already exists",
			input: model.RegisterOpts{
				Name:     "Existing User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.RegisterOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				existingUser := createTestUserWithPassword("123", "Existing User", input.Email, "oldpassword")
				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(existingUser, nil)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				require.Contains(t, err.Error(), "user have already exists")
			},
			isError: true,
		},
		{
			name: "register fails with invalid input - short name",
			input: model.RegisterOpts{
				Name:     "J",
				Email:    "john@example.com",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.RegisterOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, mongo.ErrNoDocuments)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "register fails with invalid email",
			input: model.RegisterOpts{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.RegisterOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, mongo.ErrNoDocuments)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "register fails when repository create fails",
			input: model.RegisterOpts{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.RegisterOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, mongo.ErrNoDocuments)
				repo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAuthService(&AuthServiceOpts{
				TokenManager: createTestTokenManager(),
				UserRepo:     mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.Register(context.Background(), tc.input)

			if tc.validate != nil {
				tc.validate(t, result, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type Testcase struct {
		name     string
		input    model.LoginOpts
		arrange  ArrangeFn[*AuthService, model.LoginOpts]
		validate func(t *testing.T, result *model.AuthResult, err error)
		isError  bool
	}

	testPassword := "password123"
	testUser := createTestUserWithPassword("123", "John Doe", "john@example.com", testPassword)

	testcases := []Testcase{
		{
			name: "login user successfully",
			input: model.LoginOpts{
				Email:    "john@example.com",
				Password: testPassword,
			},
			arrange: func(t *testing.T, service *AuthService, input model.LoginOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(testUser, nil)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEmpty(t, result.AccessToken)
				require.NotEmpty(t, result.RefreshToken)
				require.NotNil(t, result.User)
				require.Equal(t, "John Doe", result.User.Name())
				require.Equal(t, "john@example.com", result.User.Email())
			},
			isError: false,
		},
		{
			name: "login fails when user not found",
			input: model.LoginOpts{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			arrange: func(t *testing.T, service *AuthService, input model.LoginOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, errx.M(404, "user not found"))
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "login fails with wrong password",
			input: model.LoginOpts{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			arrange: func(t *testing.T, service *AuthService, input model.LoginOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(testUser, nil)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				require.Contains(t, err.Error(), "invalid email or password")
			},
			isError: true,
		},
		{
			name: "login fails with repository error",
			input: model.LoginOpts{
				Email:    "john@example.com",
				Password: testPassword,
			},
			arrange: func(t *testing.T, service *AuthService, input model.LoginOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByEmail(mock.Anything, input.Email).Return(nil, errors.New("database connection error"))
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAuthService(&AuthServiceOpts{
				TokenManager: createTestTokenManager(),
				UserRepo:     mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.Login(context.Background(), tc.input)

			if tc.validate != nil {
				tc.validate(t, result, err)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	type Testcase struct {
		name     string
		input    string
		arrange  ArrangeFn[*AuthService, string]
		validate func(t *testing.T, result *model.AuthResult, err error)
		isError  bool
	}

	tokenManager := createTestTokenManager()
	testUser := createTestUserWithPassword("123", "John Doe", "john@example.com", "password123")

	// Generate a valid refresh token for testing
	validRefreshToken, err := func() (string, error) {
		s := NewAuthService(&AuthServiceOpts{
			TokenManager: tokenManager,
			UserRepo:     nil, // Won't be used for token generation
		})
		return s.generateRefreshToken(testUser)
	}()
	require.NoError(t, err)

	testcases := []Testcase{
		{
			name:  "refresh token successfully",
			input: validRefreshToken,
			arrange: func(t *testing.T, service *AuthService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(testUser, nil)
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEmpty(t, result.AccessToken)
				require.Equal(t, validRefreshToken, result.RefreshToken)
				require.Nil(t, result.User) // RefreshToken doesn't return user
			},
			isError: false,
		},
		{
			name:  "refresh token fails with invalid token",
			input: "invalid.token.string",
			arrange: func(t *testing.T, service *AuthService, input string) {
				// No repo expectations since token verification fails first
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name:  "refresh token fails with expired token",
			input: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoicmVmcmVzaF90b2tlbiIsImp0aSI6IjEyMzQ1Iiwic3ViIjoiMTIzIiwiaWF0IjoxNjQwOTk1MjAwLCJleHAiOjE2NDA5OTUyMDB9.invalid",
			arrange: func(t *testing.T, service *AuthService, input string) {
				// No repo expectations since token verification fails
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name:  "refresh token fails when user not found",
			input: validRefreshToken,
			arrange: func(t *testing.T, service *AuthService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(nil, errx.M(404, "user not found"))
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name:  "refresh token fails with repository error",
			input: validRefreshToken,
			arrange: func(t *testing.T, service *AuthService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, result *model.AuthResult, err error) {
				require.Error(t, err)
				require.Nil(t, result)
			},
			isError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAuthService(&AuthServiceOpts{
				TokenManager: tokenManager,
				UserRepo:     mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.RefreshToken(context.Background(), tc.input)

			if tc.validate != nil {
				tc.validate(t, result, err)
			}
		})
	}
}

func TestAuthService_GenerateAccessToken(t *testing.T) {
	tokenManager := createTestTokenManager()
	service := NewAuthService(&AuthServiceOpts{
		TokenManager: tokenManager,
		UserRepo:     nil,
	})

	testUser := createTestUserWithPassword("123", "John Doe", "john@example.com", "password123")

	token, err := service.generateAccessToken(testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify the token can be parsed
	claims := &model.AccessTokenClaims{}
	parsedToken, err := tokenManager.VerifyTokenWithClaims(token, claims)
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	// Verify claims
	require.Equal(t, testUser.Name(), claims.Name)
	require.Equal(t, testUser.Email(), claims.Email)
	require.Equal(t, testUser.ID(), claims.Subject)
	require.Equal(t, model.AccessToken, claims.Type)
}

func TestAuthService_GenerateRefreshToken(t *testing.T) {
	tokenManager := createTestTokenManager()
	service := NewAuthService(&AuthServiceOpts{
		TokenManager: tokenManager,
		UserRepo:     nil,
	})

	testUser := createTestUserWithPassword("123", "John Doe", "john@example.com", "password123")

	token, err := service.generateRefreshToken(testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify the token can be parsed
	claims := &model.RefreshTokenClaims{}
	parsedToken, err := tokenManager.VerifyTokenWithClaims(token, claims)
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	// Verify claims
	require.Equal(t, testUser.ID(), claims.Subject)
	require.Equal(t, model.RefreshToken, claims.Type)
}
