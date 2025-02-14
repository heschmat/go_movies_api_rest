package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Check if a panic has occurred.
			if err := recover(); err != nil {
				// Setting "Connection: close" header on the response acts as a trigger
				// to make Go's HTTP server automatically close the current connection
				// after a response has been sent.
				w.Header().Set("Connection", "close")

				// N.B. The value returned by `recover()` has the type `any`.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
