package data

import (
	"database/sql"
	"errors"
)

// We'll use this in the *Get()* method when a movie could not be found.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// The *Models* struct acts as single container holding all the db models.
type Models struct {
	Movies MovieModel
}

// Initializer for the MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
