package response

type BodyProfanity struct {
	Reason        string   `json:"reason"`
	ProfaneFields []string `json:"profaneFields"`
}

func ErrProfanity(profaneKeys ...string) *Res[BodyProfanity] {
	var r Res[BodyProfanity]
	r.Status = 400
	r.Body.Reason = "profanity"
	r.Body.ProfaneFields = profaneKeys
	return &r
}
