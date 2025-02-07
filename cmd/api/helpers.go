package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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


func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	// Limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder.
	dec := json.NewDecoder(r.Body)
	// Make sure if the JSON sent by client does NOT include unkown fields.
	// As expected by the target destination.
	// This way, Decode() returns an error message like: "json: unknown field <fieldName>"
	// N.B. default behavior is to ignore.
	dec.DisallowUnknownFields()

	// Decode the body content into the target destination.
	// You have to add a *non nil pointer* as the target decode destination to `.Decode()`
	// Otherwise, you'll get `json.InvalidUnmarshalError` at runtime.
	err := dec.Decode(dst)
	if err != nil {
		// If there's an error during decoding, start the triage...
		var syntaxErr			*json.SyntaxError
		var unmarshalTypeErr	*json.UnmarshalTypeError
		var InvalidUnmarshalErr *json.InvalidUnmarshalError

		var maxBytesErr			*http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesErr):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesErr.Limit)

		// curl -d '{"title": "creed I",}' localhost:4000/v1/movie
		case errors.As(err, &syntaxErr):
			msg := "1body contains badly-formed JSON (at character %d)"
			return fmt.Errorf(msg, syntaxErr.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("2body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeErr):
			// curl -d '{"year": "10"}' localhost:4000/v1/movies
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("3abody contains incorrect JSON type for field %q", unmarshalTypeErr.Field)
			}
			// curl -d '["title", "creed I"]' localhost:4000/v1/movies
			return fmt.Errorf("3bbody contains incorrect JSON type (at character %d)", unmarshalTypeErr.Offset)

		// curl -X POST localhost:4000/v1/movies
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %q", fieldName)

		// json.InvalidUnmarshalError is returned if we pass sth that's NOT a non-nil pointer to Decode()
		// Catch this & panic.
		case errors.As(err, &InvalidUnmarshalErr):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again.
	// If the request body only contains a single JSON value, the call returns *io.EOF error*.
	// Otherwise, there's additional data in the request body, which we don't desire.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
