package types

import (
	"net/http"
	"encoding/json"
)

type ApiError struct {
	Status 	int    `json:"-"`
	Message string `json:"message"`
}

func (e *ApiError) Send(w http.ResponseWriter) {
	body, err := json.Marshal(*e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	w.Write(body)
}