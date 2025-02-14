package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new *httprouter* router instance.
	router := httprouter.New()

	// Register the relevant methods, URL patterns & handler functions for our endpoints.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)

	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(router)
}
