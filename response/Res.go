package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Res[T any] struct {
	Body   T
	Status int
}

func (r *Res[T]) Send(w http.ResponseWriter) {
	str, err := json.Marshal(r.Body)
	if err != nil {
		ErrServer(
			fmt.Sprintf("Could not marshal the following object to JSON.\n%v", r),
		).Send(w)
		return
	}
	w.WriteHeader(r.Status)
	w.Header().Add("Content-Type", "application/json")
	w.Write(str)
}

func (r *Res[T]) Print() {
	body, _ := json.Marshal(r.Body)
	fmt.Printf("\nStatus %d\nBody:\n%s\n", r.Status, body)
}


