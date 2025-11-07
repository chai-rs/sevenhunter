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

func (u *User) Validate(newUser ...bool) error {
	rules := []*v.FieldRules{
		v.Field(&u.name, v.Required, v.Length(2, 100)),
		v.Field(&u.email, v.Required, v.Length(5, 200), is.Email),
		v.Field(&u.hashedPassword, v.Required),
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
