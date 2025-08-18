package items

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
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
	if err := ViewItems(items).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ViewItem(item).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	item := Item{}
	if err := AddItem(item).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddItemPost(w http.ResponseWriter, r *http.Request) {
	item, err := getItemFromForm(r)
	if err != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	if err := h.store.Add(r.Context(), item); err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := UpdateItem(item).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateItemPost(w http.ResponseWriter, r *http.Request) {
	item, err := getItemFromForm(r)
	if err != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.store.Update(r.Context(), id, item); err != nil {
		// TODO: propagate error
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.store.Delete(r.Context(), id); err != nil {
		//TODO: propagate error
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func getItemFromForm(r *http.Request) (Item, error) {
	name := r.FormValue("name")
	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.ParseInt(quantityStr, 10, 64)
	if err != nil {
		return Item{}, err
	}
	item := Item{Name: name, Quantity: int(quantity)}
	return item, nil
}
