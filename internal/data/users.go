package data

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"password"`
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
