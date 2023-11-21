package response

import "net/http"

type BodyConflict struct {
	ConflictingFields []string `json:"conflictingFields"`
}

func ErrConflict(conflictingKeys ...string) *Res[BodyConflict] {
	var r Res[BodyConflict]
	r.Status = http.StatusConflict
	r.Body.ConflictingFields = conflictingKeys
	return &r
}
