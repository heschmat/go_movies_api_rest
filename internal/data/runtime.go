package data

import (
	"fmt"
	"strconv"
)

// Declare a custom `Runtime` type.
// This has `int32` type; but will be encoded in JSON as "r min", r being the number.
type Runtime int32

// Implement a *MarshalJSON() method* on the *Runtime*.
// This must satisfy the *json.Marshaler interface*.
// i.e., returns the expected JSON-encoded value (e.g., "39 min") as a byte slice.
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the expected format.
	jsonVal := fmt.Sprintf("%d mins", r)

	// Wrap the string `jsonVal` in double quotes before returning it.
	// Otherwise it won't be detected as a valid "JSON string".
	return []byte(strconv.Quote(jsonVal)), nil
}
