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

	return router
}
