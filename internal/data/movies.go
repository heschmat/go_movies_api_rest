package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/heschmat/go_movies_api_rest/internal/validator"
	"github.com/lib/pq"
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

type MovieModel struct {
	DB *sql.DB
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

// CRUD operations ========================================================== #
func (m MovieModel) Insert(movie *Movie) error {
	q := `INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	queryArgs := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(q, queryArgs...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	q := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(q, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Handle the errors.
	// If no matching movie found, .Scan() returns a *sql.ErrNoRows* error.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	// Return a pointer to the *Movie* struct.
	return &movie, nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	q := "DELETE FROM movies WHERE id = $1"
	result, err := m.DB.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) Update (movie *Movie) error {
	q := `UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(q, args...).Scan(&movie.Version)
}

func (m MovieModel) GetMovies(title string, genres []string) ([]*Movie, error) {
	q := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	--WHERE (LOWER(title) = LOWER($1) OR $1 = '')
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (genres @> $2 OR $2 = '{}')
	ORDER BY id`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// Execute the query.
	rows, err := m.DB.QueryContext(ctx, q, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to hold the fetched record(s).
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went ok, return the movies slice.
	return movies, nil
}
