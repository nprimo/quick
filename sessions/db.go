package sessions

import (
	"context"
	"database/sql"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(db *sql.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) Add(ctx context.Context, session Session) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", session.ID, session.UserID, session.ExpiresAt)
	return err
}

func (s *DBStore) Get(ctx context.Context, id string) (Session, error) {
	var session Session
	err := s.db.QueryRowContext(ctx, "SELECT id, user_id, expires_at FROM sessions WHERE id = ?", id).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	return session, err
}

func (s *DBStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	return err
}
