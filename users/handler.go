package users

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nprimo/quick/sessions"
	"github.com/nprimo/quick/web"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store    Store
	sessions sessions.Store
	log      *slog.Logger
}

func NewHandler(store Store, sessions sessions.Store, log *slog.Logger) Handler {
	return Handler{store: store, sessions: sessions, log: log}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) error {
	if err := Register().Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to render registration page")
	}
	return nil
}

func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return web.NewError(http.StatusBadRequest, err, "failed to parse form")
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return web.NewError(http.StatusBadRequest, errors.New("email and password are required"), "email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to hash password")
	}

	user := User{
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	if err := h.store.Add(r.Context(), user); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to create user")
	}

	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) error {
	if err := Login("").Render(r.Context(), w); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to render login page")
	}
	return nil
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return web.NewError(http.StatusBadRequest, err, "failed to parse form")
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return Login("email and password are required").Render(r.Context(), w)
	}

	user, err := h.store.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Login("wrong email or password").Render(r.Context(), w)
		}
		return web.NewError(http.StatusInternalServerError, err, "failed to get user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return Login("wrong email or password").Render(r.Context(), w)
	}

	session := sessions.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.sessions.Add(r.Context(), session); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to create session")
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/items", http.StatusFound)
	return nil
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}

	if err := h.sessions.Delete(r.Context(), cookie.Value); err != nil {
		return web.NewError(http.StatusInternalServerError, err, "failed to delete session")
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}
