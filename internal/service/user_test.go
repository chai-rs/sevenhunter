package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chai-rs/sevenhunter/internal/model"
	"github.com/chai-rs/sevenhunter/internal/model/mocks"
	errx "github.com/chai-rs/sevenhunter/pkg/error"
	. "github.com/chai-rs/sevenhunter/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestUser(id, name, email string) *model.User {
	user, err := model.NewUser(model.UserOpts{
		ID:             id,
		Name:           name,
		Email:          email,
		HashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		panic("Failed to create test user: " + err.Error())
	}
	return user
}

func TestUserService_Count(t *testing.T) {
	type Testcase struct {
		name     string
		arrange  ArrangeFn[*UserService, any]
		expected int64
		isError  bool
	}

	testcases := []Testcase{
		{
			name: "count users successfully",
			arrange: func(t *testing.T, service *UserService, _ any) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Count(mock.Anything).Return(int64(5), nil)
			},
			expected: 5,
			isError:  false,
		},
		{
			name: "count returns zero when no users",
			arrange: func(t *testing.T, service *UserService, _ any) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Count(mock.Anything).Return(int64(0), nil)
			},
			expected: 0,
			isError:  false,
		},
		{
			name: "count fails with repository error",
			arrange: func(t *testing.T, service *UserService, _ any) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Count(mock.Anything).Return(int64(0), errors.New("database error"))
			},
			expected: 0,
			isError:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserService(UserServiceOpts{
				UserRepo: mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, nil)
			}

			result, err := s.Count(context.Background())
			if tc.isError {
				require.Error(t, err)
				require.Equal(t, tc.expected, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUserService_List(t *testing.T) {
	user1 := createTestUser("1", "John Doe", "john@example.com")
	user2 := createTestUser("2", "Jane Smith", "jane@example.com")

	type Testcase struct {
		name     string
		input    model.ListUserOpts
		arrange  ArrangeFn[*UserService, model.ListUserOpts]
		expected []model.User
		isError  bool
	}

	testcases := []Testcase{
		{
			name:  "list users successfully with empty result",
			input: model.ListUserOpts{},
			arrange: func(t *testing.T, service *UserService, input model.ListUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().List(mock.Anything, input).Return([]model.User{}, nil)
			},
			expected: []model.User{},
			isError:  false,
		},
		{
			name: "list users successfully with results",
			input: model.ListUserOpts{
				Limit:   10,
				SortAsc: true,
			},
			arrange: func(t *testing.T, service *UserService, input model.ListUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().List(mock.Anything, input).Return([]model.User{*user1, *user2}, nil)
			},
			expected: []model.User{*user1, *user2},
			isError:  false,
		},
		{
			name: "list users with pagination cursor",
			input: model.ListUserOpts{
				Cursor:  "cursor123",
				Limit:   5,
				SortAsc: false,
			},
			arrange: func(t *testing.T, service *UserService, input model.ListUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().List(mock.Anything, input).Return([]model.User{*user1}, nil)
			},
			expected: []model.User{*user1},
			isError:  false,
		},
		{
			name:  "list users fails with repository error",
			input: model.ListUserOpts{},
			arrange: func(t *testing.T, service *UserService, input model.ListUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().List(mock.Anything, input).Return(nil, errors.New("database connection error"))
			},
			expected: nil,
			isError:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserService(UserServiceOpts{
				UserRepo: mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.List(context.Background(), tc.input)
			if tc.isError {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUserService_Get(t *testing.T) {
	user := createTestUser("123", "John Doe", "john@example.com")

	type Testcase struct {
		name     string
		input    string
		arrange  ArrangeFn[*UserService, string]
		expected *model.User
		isError  bool
	}

	testcases := []Testcase{
		{
			name:  "get user successfully",
			input: "123",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input).Return(user, nil)
			},
			expected: user,
			isError:  false,
		},
		{
			name:  "get user fails when user not found",
			input: "nonexistent",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input).Return(nil, errx.M(404, "user not found"))
			},
			expected: nil,
			isError:  true,
		},
		{
			name:  "get user fails with repository error",
			input: "123",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input).Return(nil, errors.New("database error"))
			},
			expected: nil,
			isError:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserService(UserServiceOpts{
				UserRepo: mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.Get(context.Background(), tc.input)
			if tc.isError {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	existingUser := createTestUser("123", "John Doe", "john@example.com")

	type Testcase struct {
		name     string
		input    model.UpdateUserOpts
		arrange  ArrangeFn[*UserService, model.UpdateUserOpts]
		validate func(t *testing.T, result *model.User)
		isError  bool
	}

	testcases := []Testcase{
		{
			name: "update user successfully",
			input: model.UpdateUserOpts{
				ID:    "123",
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			arrange: func(t *testing.T, service *UserService, input model.UpdateUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input.ID).Return(existingUser, nil)
				repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Name() == "John Updated" && u.Email() == "john.updated@example.com"
				})).Return(nil)
			},
			validate: func(t *testing.T, result *model.User) {
				require.NotNil(t, result)
				require.Equal(t, "John Updated", result.Name())
				require.Equal(t, "john.updated@example.com", result.Email())
			},
			isError: false,
		},
		{
			name: "update user fails when user not found",
			input: model.UpdateUserOpts{
				ID:    "nonexistent",
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			arrange: func(t *testing.T, service *UserService, input model.UpdateUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input.ID).Return(nil, errx.M(404, "user not found"))
			},
			validate: func(t *testing.T, result *model.User) {
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "update user fails with invalid email",
			input: model.UpdateUserOpts{
				ID:    "123",
				Name:  "John Updated",
				Email: "invalid-email",
			},
			arrange: func(t *testing.T, service *UserService, input model.UpdateUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input.ID).Return(existingUser, nil)
				// Update will not be called because validation fails
			},
			validate: func(t *testing.T, result *model.User) {
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "update user fails with empty name",
			input: model.UpdateUserOpts{
				ID:    "123",
				Name:  "",
				Email: "john@example.com",
			},
			arrange: func(t *testing.T, service *UserService, input model.UpdateUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input.ID).Return(existingUser, nil)
			},
			validate: func(t *testing.T, result *model.User) {
				require.Nil(t, result)
			},
			isError: true,
		},
		{
			name: "update user fails when repository update fails",
			input: model.UpdateUserOpts{
				ID:    "123",
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			arrange: func(t *testing.T, service *UserService, input model.UpdateUserOpts) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().FindByID(mock.Anything, input.ID).Return(existingUser, nil)
				repo.EXPECT().Update(mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			validate: func(t *testing.T, result *model.User) {
				require.Nil(t, result)
			},
			isError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserService(UserServiceOpts{
				UserRepo: mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			result, err := s.Update(context.Background(), tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.validate != nil {
				tc.validate(t, result)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type Testcase struct {
		name    string
		input   string
		arrange ArrangeFn[*UserService, string]
		isError bool
	}

	testcases := []Testcase{
		{
			name:  "delete user successfully",
			input: "123",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Delete(mock.Anything, input).Return(nil)
			},
			isError: false,
		},
		{
			name:  "delete user fails when user not found",
			input: "nonexistent",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Delete(mock.Anything, input).Return(errx.M(404, "user not found"))
			},
			isError: true,
		},
		{
			name:  "delete user fails with repository error",
			input: "123",
			arrange: func(t *testing.T, service *UserService, input string) {
				repo, ok := service.userRepo.(*mocks.MockUserRepo)
				require.True(t, ok)

				repo.EXPECT().Delete(mock.Anything, input).Return(errors.New("database error"))
			},
			isError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserService(UserServiceOpts{
				UserRepo: mocks.NewMockUserRepo(t),
			})

			if tc.arrange != nil {
				tc.arrange(t, s, tc.input)
			}

			err := s.Delete(context.Background(), tc.input)
			if tc.isError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
