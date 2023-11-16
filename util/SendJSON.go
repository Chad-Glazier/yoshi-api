package util

import (
	"encoding/json"
	"net/http"
)

// Sends the given object as a JSON body. If the object fails to marshal to
// JSON for some reason, it will send a `500` error instead.
func SendJSON(w http.ResponseWriter, v any) {
	str, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(str)
}

// Sends the given string as a JSON object. I.e., it sets the Content-Type 
// header and writes the string to the response.
func SendAsJSON(w http.ResponseWriter, s string) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(s))
} 
