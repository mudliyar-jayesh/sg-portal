package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// ParseJSONBody is a generic function to read the request body and decode it into the provided type `T`.
func ParseJSONBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}

// RespondJSON is a generic function to send any type of response as JSON with the provided status code.
func RespondJSON[T any](w http.ResponseWriter, statusCode int, payload *T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// ParseUintParam extracts and parses an unsigned integer from query parameters
func ParseUintParam(r *http.Request, param string) (uint64, error) {
	values := r.URL.Query()
	paramStr := values.Get(param)
	if paramStr == "" {
		return 0, errors.New("missing parameter: " + param)
	}

	value, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid parameter: " + param)
	}
	return value, nil
}

// HandleError is used to send an HTTP error response with a custom message and status code.
func HandleError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
