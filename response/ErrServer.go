package response

import "net/http"

type BodyServer struct {
	Message string `json:"message"`
}

func ErrServer(msg string) *Res[BodyServer] {
	var r Res[BodyServer]
	r.Status = http.StatusInternalServerError
	r.Body = BodyServer{
		Message: msg,
	}
	return &r
}

func ErrServerUnkown() *Res[BodyServer] {
	return ErrServer(
		"unknown internal server error",
	)
}
