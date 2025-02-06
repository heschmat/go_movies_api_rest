package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Define an envelope type (for json formatting).
type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	// .ParamsFromContext() returns a slice containing names & values of the parameters.
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// w: the destination http.ResponseWriter
// data: data to encode to JSON
// status: the HTTP status code to send
// headers: header map containing additional HTTP headers to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, data envelope, status int, headers http.Header) error {
	// json.Marshal() returns a `[]byte` slice containing the encoded JSON.
	// js, err := json.Marshal(data)
	// *MarshalIndent* for better readability of JSON in terminal.
	js, err := json.MarshalIndent(data, "", "\t")  // no prefix
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view the JSON data in terminal.
	js = append(js, '\n')

	// Loop through the header map & add each header to the http.ResponseWriter header map.
	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
