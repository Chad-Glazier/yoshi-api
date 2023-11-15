package mw

import (
	"net/http"
)

// Used to enforce a method constraint. Note that this should never be used
// before the `mw.Cors` middleware.
//
// May send back the following errors:
//	- `http.StatusMethodNotAllowed` (no body)
func Method(desiredMethod string) MiddlewareFunc {
	return func(res *Resources, w http.ResponseWriter, r *http.Request) (bool, CleanupFunc) {
		if r.Method != desiredMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return false, nil
		}
		return true, nil
	}
}
