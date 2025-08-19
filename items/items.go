package items

import (
	"context"
	"time"
)

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store interface {
	Add(ctx context.Context, item Item) error
	All(ctx context.Context) ([]Item, error)
	Get(ctx context.Context, id int) (Item, error)
	Update(ctx context.Context, id int, item Item) error
	Delete(ctx context.Context, id int) error
}
