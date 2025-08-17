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

	mux.HandleFunc("GET /items/{id}/update", itemHandler.UpdateItem)
	mux.HandleFunc("POST /items/{id}", itemHandler.UpdateItemPost)

	mux.HandleFunc("GET /items/{id}/delete", itemHandler.DeleteItem)

	mux.HandleFunc("GET /items/new", itemHandler.AddItem)
	mux.HandleFunc("POST /items/new", itemHandler.AddItemPost)

	wrapped := LoggerMiddleware(mux, log)
	return wrapped
}
