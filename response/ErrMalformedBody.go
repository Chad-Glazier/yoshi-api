package response

import (
	"net/http"
	"encoding/json"
)

type Exemplar[T any] interface {
	Example() T
}

type BodyMalformedBody[T Exemplar[T]] struct {
	Reason   string `json:"reason"`
	Expected T      `json:"expected"`
	Received any    `json:"received"`
}

func ErrMalformedBody[T Exemplar[T]](r *http.Request) *Res[BodyMalformedBody[T]] {
	var expected T
	expected = expected.Example()
	
	var received any
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&received)

	res := Res[BodyMalformedBody[T]]{
		Status: http.StatusBadRequest,
		Body: BodyMalformedBody[T]{
			Reason: "malformed request body",
			Expected: expected,
			Received: received,
		},
	}

	return &res
}