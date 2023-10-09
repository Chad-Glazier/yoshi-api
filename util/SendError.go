package util

import (
	"encoding/json"
	"net/http"
	"yoshi-api/types"
)

func SendError(w http.ResponseWriter, errorInfo types.ErrorInfo) {
	body, err := json.Marshal(types.ApiError{
		Message: errorInfo.Message,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorInfo.Status)
	w.Write(body)
}
