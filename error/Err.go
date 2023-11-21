package error

import (
	"encoding/json"
	"log"
	"net/http"
)

type Err[T any] struct {
	Body 	T
	Status 	int
}

type ServerBody struct {
	Message string `json:"message"`
}

func Server(msg string, status ...int) *Err[ServerBody] {
	var e Err[ServerBody]
	e.Status = 500
	e.Body = ServerBody{
		Message: msg,
	}
	return &e
}

func SendErr(w http.ResponseWriter, e Err[any]) {
	str, err := json.Marshal(e.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("could not marshal the following object to JSON\n%v", e)
		return
	}
	w.WriteHeader(e.Status)
	w.Header().Add("Content-Type", "application/json")
	w.Write(str)
}
