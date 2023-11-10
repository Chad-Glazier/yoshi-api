package util

import (
	"net/http"
)

// sets the `session_id` cookie on a request, with the appropriate settings.
func SetSessionCookie(sessionId string, w http.ResponseWriter) {
	sessionCookie := &http.Cookie{
		Name:  "session_id",
		Value: sessionId,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, sessionCookie)
}