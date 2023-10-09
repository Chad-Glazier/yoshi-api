package middleware

import (
	"net/http"

	"yoshi-api/types"
)

func Authenticated(r *http.Request) (err *types.ErrorInfo) {
	authenticated := false
	// authentication logic here!
	if !authenticated {
		return &types.ErrorInfo{
			Status:  http.StatusUnauthorized,
			Message: "Authorization failed",
		}
	}
	return nil // this means "continue with the pipeline of middleware/handlers"
}
