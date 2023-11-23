package response

import (
	"net/http"
)

func OKEmpty() *Res[struct{}] {
	var r Res[struct{}]
	r.Status = http.StatusOK
	return &r
}

func OK[T any](body T) *Res[T] {
	var r Res[T]
	r.Body = body
	r.Status = http.StatusOK
	return &r
}
