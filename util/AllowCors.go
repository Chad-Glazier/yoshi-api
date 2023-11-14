package util

import (
	"net/http"
)

// sets the appropriate headers to allow for cross-origin requests. Browsers might
// also make preflight requests, in which case this function will also send an appropriate
// response and return `true` to indicate that you shouldn't proceed with the request.
//
// TL/DR: Ignore the request if this returns `true`.
func AllowCors(w http.ResponseWriter, r *http.Request) (stop bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,PATCH,DELETE,POST,PUT")
	w.Header().Set(
		"Access-Control-Allow-Headers", 
		"X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version",
	)

	// Handle preflight requests (OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}
