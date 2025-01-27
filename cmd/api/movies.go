package main

import (
	"fmt"
	"net/http"
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

	// Otherwise, for now, interpolate the movie ID in a placeholder response.
	fmt.Fprintf(w, "Detail of movie with id: %d\n", id)
}
