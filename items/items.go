package items

import (
	"context"
	"time"
)

type Item struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store interface {
	Add(ctx context.Context, item Item, userID int) error
	All(ctx context.Context, userID int) ([]Item, error)
	Get(ctx context.Context, id int, userID int) (Item, error)
	Update(ctx context.Context, id int, item Item, userID int) error
	Delete(ctx context.Context, id int, userID int) error
}
