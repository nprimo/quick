package users

import (
	"context"
	"database/sql"
)

type DBStore struct {
	db *sql.DB
}

var _ Store = (*DBStore)(nil)

func NewDBStore(db *sql.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) Add(ctx context.Context, user User) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (email, hashed_password) VALUES (?, ?)", user.Email, user.HashedPassword)
	return err
}

func (s *DBStore) Get(ctx context.Context, id int) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, "SELECT id, email, hashed_password, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (s *DBStore) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, "SELECT id, email, hashed_password, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}
