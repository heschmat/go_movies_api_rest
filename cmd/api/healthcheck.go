package main

import (
	"net/http"
)

// // Writes a plain-text response with information about:
// // the application status, operating environment & version.
// func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintln(w, "status: available")
// 	// fmt.Fprintf(w, "environment: %s\n", app.config.env)
// 	// fmt.Fprintf(w, "version: %s\n", version)

// 	// %q wraps the interpolated values in double-quotes.
// 	js := `{"status": "available", "environment": %q, "version": %q}`
// 	js = fmt.Sprintf(js, app.config.env, version)

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write([]byte(js))
// }


func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an *envelop map*, containing the info for the response.
	info := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, info, http.StatusOK, nil)
	if err != nil {
		app.logger.Error(err.Error())
		msg := "Server encountered an issue & could not process your request"
		http.Error(w, msg, http.StatusInternalServerError)
	}
}
