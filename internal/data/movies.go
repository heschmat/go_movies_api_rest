package data

import (
	"time"

	"github.com/heschmat/go_movies_api_rest/internal/validator"
)

// struct tags control how the keys appear in the JSON-encoded output.
// N.B. The `-` directive removes the struct field from being encoded as JSON data.
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // TimeStamp for when the movie is added to our db
	Title     string    `json:"title"`
	Year      int32     `json:"year"`    // Movie release year
	Runtime   Runtime     `json:"runtime"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres"`  // Slice of genres for the movie (drama, comedy, romance ...)
	Version   int32     `json:"version"` // starts at 1; will be incremented each time movie info is updated
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
}
