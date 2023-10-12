package types

import (
	"net/http"
)

type MiddlewareFunc func(r *http.Request) (err *ApiError)
