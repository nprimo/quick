package web

import (
	"log/slog"
	"net/http"

	"github.com/nprimo/quick/items"
)

func Router(
	itemHandler items.Handler,
	log *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(`ciao mamma`))
	})

	mux.HandleFunc("GET /items", itemHandler.ListItems)
	mux.HandleFunc("GET /items/{id}", itemHandler.GetItem)

	wrapped := LoggerMiddleware(mux, log)
	return wrapped
}
