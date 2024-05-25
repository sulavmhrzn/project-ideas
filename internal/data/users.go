package data

import (
	"time"

	"github.com/sulavmhrzn/projectideas/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type password struct {
	PlainPassword  string
	HashedPassword []byte
}

func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.PlainPassword = plainPassword
	p.HashedPassword = hash
	return nil
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 10, "username", "must not be greater than 10 characters long")
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.ValidEmail(user.Email), "email", "must be a valid email address")
	v.Check(user.Password.PlainPassword != "", "password", "must be provided")
	v.Check(len(user.Password.PlainPassword) < 72, "password", "must not be greater than 72 characters long")
	v.Check(len(user.Password.PlainPassword) > 10, "password", "must be greater than 10 characters long")
}
