package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/heschmat/go_movies_api_rest/internal/data"
)

// corresponding endpoint: "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "A new movie added.")
}

// corresponding endpoint: "GET /v1/movies/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Create a new instance of the Movie struct
	// ID is extracted for the URL, the rest are dummy data for now.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Creed I",
		Year:      2015,
		Runtime:   133,
		Genres:    []string{"drama", "action", "boxing", "sport"},
		Version:   1,
	}

	err = app.writeJSON(w, movie, http.StatusOK, nil)
	if err != nil {
		app.logger.Error(err.Error())
		msg := "Server encountered an issue & could not process your request"
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
