package web

import (
	"log/slog"
	"net/http"
)

type Error struct {
	Code          int
	Message       string
	InternalError error
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, err error, message string) Error {
	if len(message) == 0 {
		message = http.StatusText(code)
	}
	return Error{
		Code:          code,
		InternalError: err,
		Message:       message,
	}
}

// HandlerFuncWithErr is a custom handler function that returns an error.
// This allows us to centralize error handling.
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler is a middleware that handles errors from HandlerFuncWithErr.
func ErrorHandler(log *slog.Logger) func(HandlerFuncWithErr) http.Handler {
	return func(handler HandlerFuncWithErr) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := handler(w, r)
			if err != nil {
				err, ok := err.(Error)
				if !ok {
					log.Error("unknown error",
						"err", err)
				}
				log.Error("ERROR:",
					"err", err.InternalError)
			}
		})
	}
}
