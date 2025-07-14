package items

import (
	"net/http"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) Handler {
	return Handler{store: store}
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ViewItems(items).Render(r.Context(), w)
}

