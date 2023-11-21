package response

import "net/http"

type BodyUnauthorized struct {
	Message string	`json:"message"`
}

func ErrUnauthorized() *Res[BodyUnauthorized] {
	var r Res[BodyUnauthorized]
	r.Status = http.StatusUnauthorized
	r.Body = BodyUnauthorized{
		Message: "request did not contain a valid session_id",
	}
	return &r
}
