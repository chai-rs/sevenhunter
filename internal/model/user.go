package model

import (
	"net/http"
	"time"

	errx "github.com/chai-rs/sevenhunter/pkg/error"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id             string
	name           string
	email          string
	hashedPassword string
	createdAt      time.Time
}

type UserOpts struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

func NewUser(opts UserOpts) (*User, error) {
	u := &User{
		id:             opts.ID,
		name:           opts.Name,
		email:          opts.Email,
		hashedPassword: opts.HashedPassword,
		createdAt:      opts.CreatedAt,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate(newUser ...bool) error {
	rules := []*v.FieldRules{
		v.Field(&u.name, v.Required, v.Length(2, 100)),
		v.Field(&u.email, v.Required, v.Length(5, 200), is.Email),
		v.Field(&u.hashedPassword, v.Required),
		v.Field(&u.createdAt, v.Required),
	}

	// not new user, validate the auto-generated fields
	if len(newUser) == 0 || !newUser[0] {
		rules = append(rules, v.Field(&u.hashedPassword, v.Length(60, 60)))
	}

	return v.ValidateStruct(u, rules...)
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

func (u *User) ComparePassword(password string) error {
	if bcrypt.CompareHashAndPassword([]byte(u.hashedPassword), []byte(password)) != nil {
		return errx.M(http.StatusBadRequest, "invalid email or password")
	}
	return nil
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

type UpdateUserOpts struct {
	ID    string
	Name  string
	Email string
}

func (s *User) Update(opts UpdateUserOpts) error {
	s.email = opts.Email
	s.name = opts.Name

	if err := s.Validate(); err != nil {
		return errx.E(http.StatusBadRequest, err)
	}

	return nil
}

type CreateUserOpts struct {
	Name     string
	Email    string
	Password string
}

func NewCreateUser(opts CreateUserOpts) (*User, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(opts.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := User{
		name:           opts.Name,
		email:          opts.Email,
		hashedPassword: string(hashedPasswordBytes),
		createdAt:      time.Now(),
	}

	if err := u.Validate(true); err != nil {
		return nil, errx.E(http.StatusBadRequest, err)
	}

	return &u, nil
}

type ListUserOpts struct {
	Cursor  string
	Limit   int
	SortAsc bool
}

func (opts *ListUserOpts) GetLimit() int {
	if opts.Limit <= 0 {
		return 10
	}
	if opts.Limit > 100 {
		return 100
	}
	return opts.Limit
}
