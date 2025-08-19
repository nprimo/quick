package sessions

import (
	"context"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CSRFToken string    `json:"csrf_token"`
}

type Store interface {
	Add(ctx context.Context, session Session) error
	Get(ctx context.Context, id string) (Session, error)
	Save(ctx context.Context, session Session) error
	Delete(ctx context.Context, id string) error
}
