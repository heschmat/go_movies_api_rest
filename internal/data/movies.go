package data

import "time"

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
