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

// HandlerFuncWithErr is a custom handler function that returns an error.
// This allows us to centralize error handling.
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) error

// ErrorMux is a custom ServeMux that wraps handlers to manage errors.
// It embeds http.ServerMux to inherit its functionality.
type ErrorMux struct {
	*http.ServeMux
	log *slog.Logger
}

// NewErrorMux creates and returns a new ErrorMux.
func NewErrorMux(log *slog.Logger) *ErrorMux {
	return &ErrorMux{
		ServeMux: http.NewServeMux(),
		log:      log,
	}
}

// HandleErrorFunc registers a custom ErrorHandler for the given pattern.
// It wraps the ErrorHandler in a standard http.HandlerFunc to check for errors.
func (mux *ErrorMux) HandleErrorFunc(pattern string, handler HandlerFuncWithErr) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			err, ok := err.(Error)
			if !ok {
				mux.log.Error("unknown error",
					"err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			mux.log.Error("ERROR:",
				"err", err)
			http.Error(w, http.StatusText(err.Code), err.Code)
			return
		}
	}
	mux.HandleFunc(pattern, wrappedHandler)
}
