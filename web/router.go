package web

import (
	"net/http"

	"github.com/nprimo/quick/items"
)

func Router(itemHandler items.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ciao mamma`))
	})

	mux.HandleFunc("/items", itemHandler.ListItems)

	return mux
}