package mw

import (
	"net/http"
	"os"
)

// Configures CORS, and handles `OPTION` requests.
func Cors(res *Resources, w http.ResponseWriter, r *http.Request) (bool, CleanupFunc){
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_ORIGIN"))
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET,OPTIONS,PATCH,DELETE,POST,PUT",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Credentials, X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version",
	)

	if (r.Method == http.MethodOptions) {
		w.WriteHeader(http.StatusOK)
		return false, nil
	}
	return true, nil
}
