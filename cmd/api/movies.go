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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id , err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"message": "movie successfully deleted"}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
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

	// Hold the expected data from the client.
	var input struct {
		Title	*string		`json:"title"`
		Year	*int32		`json:"year"`
		Runtime	*data.Runtime	`json:"runtime"`
		Genres	[]string	`json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the appropriate fields of the movie record.
	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *&movie.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Hold the expected values from the request query string.
	var input struct {
		Title		string
		Genres		[]string
		Page		int
		PageSize	int
		Sort		string
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSVString(qs, "genres", []string{})

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 10, v)

	// By default the sort is ascending by id.
	input.Sort = app.readString(qs, "sort", "id")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, err := app.models.Movies.GetMovies(input.Title, input.Genres)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"movies": movies}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
