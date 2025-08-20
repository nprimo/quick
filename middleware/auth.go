package middleware

import (
	"net/http"

	"github.com/nprimo/quick/sessions"
)

func RequireLogin(next http.Handler) (http.Handler, error) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if userID := sessions.GetUserID(r.Context()); userID == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn), nil
}
