package main

import (
	"fmt"
	"net/http"
)

// a generic helper for logging an error message
func (app *application) logError(r *http.Request, err error) {
	// Also log the current request method & URL.
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// a generic helper for sending JSON-formatted error messages to the client.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, envelope{"error": message}, status, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// When application encounters an unexpected problem at runtime.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "Server encountered an issue & could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "The requested resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("The %q method is NOT supported for this resource.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
