package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// JSON encodes to JSON format given data and sends a response
// to the client with a given http code and encoded data.
func JSON(w http.ResponseWriter, code int, data interface{}) {
	obj, err := json.Marshal(data)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(obj)
}

// readJSON decodes request body to the given destination(usually model struct).
// Returns an error on failure.
func ReadJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	// Create a new decoder and check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {
		// If an error occurred, send an error mapped to JSON decoding error
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"request body contains badly-formatted JSON (at character %d)",
				syntaxError.Offset,
			)
		// Type error
		case errors.As(err, &unmarshalTypeError):
			// If there's an info for struct field, show what field contains an error
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"request body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}

			return fmt.Errorf(
				"request body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)
		// Unmarshall error
		case errors.As(err, &invalidUnmarshalError):
			// We are panicing here because this is unexpected error
			panic(err)
		// Empty JSON error
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		// Unknown field error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown key %s", fieldName)

		// Return error as-is
		default:
			return err
		}
	}

	// Decode one more time to check wheter here is another JSON object
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must only contain single JSON value")
	}

	return nil
}
