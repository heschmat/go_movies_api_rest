package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/heschmat/go_movies_api_rest/internal/data"
	"github.com/heschmat/go_movies_api_rest/internal/validator"
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

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new validator instance.
	v := validator.New()

	// Copy the vaulues from the input struct into a new *Movie* struct.
	movie := &data.Movie{
		Title: 		input.Title,
		Year: 		input.Year,
		Runtime: 	data.Runtime(input.Runtime),
		Genres: 	input.Genres,
	}
	// If any of the checks failed, send `422 unprocessable entity` error.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	// Let the client know which URL the newly-created resource can be found at.
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, envelope{"movie": movie}, http.StatusCreated, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// corresponding endpoint: "GET /v1/movies/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
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
