package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var Validate = validator.New() // singleton instance of the validator

// ParseJSON parses the request body as JSON
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload) // decode the request body to the payload

}

// WriteJSON writes the value v to the response writer w as JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json") // set the content type to JSON
	w.WriteHeader(status) // set the status code (e.g. 200 OK, 404 Not Found, etc.)
	return json.NewEncoder(w).Encode(v) // encode the value to JSON and write it to the response writer
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()}) // write the error message as JSON
}

// GetTokenFromRequest extracts the token from the request
func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")
	
	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}