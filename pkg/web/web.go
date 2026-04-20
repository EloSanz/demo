package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Error is an HTTP error with a status code and user-facing message.
type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func NewError(statusCode int, message string) error {
	return &Error{StatusCode: statusCode, Message: message}
}

func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return &Error{StatusCode: http.StatusBadRequest, Message: "malformed request body: " + err.Error()}
	}
	return nil
}

func EncodeJSON(w http.ResponseWriter, payload any, statusCode int) error {
	if payload == nil {
		w.WriteHeader(statusCode)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("encoding response: %w", err)
	}
	return nil
}

// Param extracts a named path parameter using Go 1.22 ServeMux {name} syntax.
func Param(r *http.Request, name string) (string, error) {
	value := r.PathValue(name)
	if value == "" {
		return "", &Error{StatusCode: http.StatusBadRequest, Message: "missing path parameter: " + name}
	}
	return value, nil
}

// ParamInt extracts a named path parameter and converts it to int64.
func ParamInt(r *http.Request, name string) (int64, error) {
	raw, err := Param(r, name)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, &Error{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("path parameter %q must be an integer", name)}
	}
	return v, nil
}

// WriteError writes a JSON error response. Uses *Error status/message; defaults to 500.
func WriteError(w http.ResponseWriter, err error) {
	var webErr *Error
	if errors.As(err, &webErr) {
		http.Error(w, `{"error":"`+webErr.Message+`"}`, webErr.StatusCode)
		return
	}
	http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
}
