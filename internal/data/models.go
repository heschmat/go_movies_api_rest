package data

import "database/sql"

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
