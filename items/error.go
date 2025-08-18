package items

import "net/http"

type Error struct {
	Internal   error `json:"-"`
	Message    any   `json:"error"`
	StatusCode int   `json:"-"`
}

func (e Error) Error() string {
	return http.StatusText(e.StatusCode)
}
