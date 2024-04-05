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

// sendJSON sends a success code and json encoded content
func sendJSON(w http.ResponseWriter, code int, content any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(content)
}

// sendHTML sends a success code of 200 and the HTML content provided
func sendHTML(w http.ResponseWriter, html []byte) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(html)
}

func deleteElement[T any](slice []T, index int) []T {
	// Copy the elements following the index one position to the left.
	copy(slice[index:], slice[index+1:])
	// Return the slice without the last element.
	return slice[:len(slice)-1]
}
