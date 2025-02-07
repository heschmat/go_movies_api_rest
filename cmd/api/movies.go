package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/heschmat/go_movies_api_rest/internal/data"
)

// corresponding endpoint: "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct
	// to hold the information that we expect to be in the HTTP request body.
	var input struct {
		Title   string		`json:"title"`
		Year	int32		`json:"year"`
		Runtime int32		`json:"runtime"`
		Genres  []string	`json:"genres"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Dump the contents of the input struct in a HTTP response.
	// `%+v` prints struct fields along with their names. {"title": "still alice", "year": 2014}
	// `%v` prints the value in its default format: {"still alice", 2014}
	fmt.Fprintf(w, "%+v\n", input)
}

// corresponding endpoint: "GET /v1/movies/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFoundResponse(w, r)
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

	// Pass an *envelop map* instead of passing the plain movie struct.
	err = app.writeJSON(w, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		// app.logger.Error(err.Error())
		// msg := "Server encountered an issue & could not process your request"
		// http.Error(w, msg, http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}
