package items

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/web"
)

type Handler struct {
	store Store
	log   *slog.Logger
}

func NewHandler(store Store, log *slog.Logger) Handler {
	return Handler{store: store, log: log}
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	items, err := h.store.All(r.Context(), userID)
	if err != nil {
		return web.NewError(http.StatusBadRequest, err, "")
	}
	if err := ViewItems(items).Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return web.NewError(http.StatusBadRequest, err, "")
	}
	item, err := h.store.Get(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewError(http.StatusNotFound, err, "")
		}
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	if err := ViewItem(item).Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	item := Item{}
	if err := AddItem(item).Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}

func (h *Handler) AddItemPost(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	item, err := getItemFromForm(r)
	if err != nil {
		return web.NewError(http.StatusBadRequest, err, "")
	}
	if err := h.store.Add(r.Context(), item, userID); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	http.Redirect(w, r, "/items", http.StatusFound)
	return nil
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	item, err := h.store.Get(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewError(http.StatusNotFound, err, "")
		}
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	if err := UpdateItem(item).Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	return nil
}

func (h *Handler) UpdateItemPost(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	item, err := getItemFromForm(r)
	if err != nil {
		return web.NewError(http.StatusBadRequest, err, "")
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return web.NewError(http.StatusBadRequest, err, "")
	}
	if err := h.store.Update(r.Context(), id, item, userID); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	http.Redirect(w, r, "/items", http.StatusFound)
	return nil
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) error {
	userID := sessions.GetUserID(r.Context())
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	if err := h.store.Delete(r.Context(), id, userID); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "")
	}
	http.Redirect(w, r, "/items", http.StatusFound)
	return nil
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
