package users

import (
	"context"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Store interface {
	Add(ctx context.Context, user User) error
	Get(ctx context.Context, id int) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
