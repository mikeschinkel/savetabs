package restapi

import (
	"encoding/json"
	"net/http"
)

func ptr[T any](a T) *T {
	return &a
}

// sendError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendError(w http.ResponseWriter, code int, message string) {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(petErr)
}

// sendSuccess sends a success code and json encoded content
func sendSuccess(w http.ResponseWriter, code int, content any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(content)
}
