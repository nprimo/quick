package middleware

import (
	"net/http"

	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/web"
)

func RequireLogin(next web.HandlerFuncWithError) web.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		if userID := sessions.GetUserID(r.Context()); userID == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}
		return next(w, r)
	}
}
