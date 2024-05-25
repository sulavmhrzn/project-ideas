package data

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
)

type Model struct {
	User UserModel
}

func NewModel(db *sql.DB) Model {
	return Model{
		User: UserModel{DB: db},
	}
}
