package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	UserId    int       `json:"-"`
	Token     string    `json:"token"`
	Scope     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenModel struct {
	DB *sql.DB
}

func generateToken(userId int, ttl time.Duration, scope string) (*Token, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token := &Token{
		UserId:    userId,
		Token:     base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes),
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl),
	}
	return token, nil
}

func (m *TokenModel) New(userId int, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s", token.Token)
	err = m.Insert(token)
	return token, err
}

func (m *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens 
	(userId, token, scope, expires_at)
	VALUES 
	($1, $2, $3, $4)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []any{token.UserId, token.Token, token.Scope, token.ExpiresAt}
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
