package data

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrNoRows            = errors.New("no rows found")
)

type Model struct {
	User  UserModel
	Token TokenModel
	Idea  IdeaModel
}

func NewModel(db *sql.DB) Model {
	return Model{
		User:  UserModel{DB: db},
		Token: TokenModel{DB: db},
		Idea:  IdeaModel{DB: db},
	}
}
