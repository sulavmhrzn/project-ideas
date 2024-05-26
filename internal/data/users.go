package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sulavmhrzn/projectideas/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) IsAnonymousUser() bool {
	return u == AnonymousUser
}

type UserModel struct {
	DB *sql.DB
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

func (p *password) Compare(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(plainPassword))
	return err == nil
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

func (m UserModel) Insert(user *User) (*User, error) {
	query := `INSERT INTO users 
	(username, email, hash_password)
	VALUES 
	($1, $2, $3)
	RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []any{user.Username, user.Email, user.Password.HashedPassword}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return nil, ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}
	return user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, hash_password
	FROM users
	WHERE email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password.HashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetForToken(token string) (*User, error) {
	query := `
	SELECT id, username, email, created_at 
	FROM users
	JOIN tokens
	ON tokens.userId = users.id
	WHERE tokens.token = $1
	AND expires_at > now()`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	err := m.DB.QueryRowContext(ctx, query, token).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &user, nil
}
