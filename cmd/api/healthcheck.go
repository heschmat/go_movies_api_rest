package main

import (
	"fmt"
	"net/http"
)

// Writes a plain-text response with information about:
// the application status, operating environment & version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
