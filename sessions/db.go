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
	_, err := s.db.ExecContext(ctx, "INSERT INTO sessions (id, user_id, expires_at, csrf_token) VALUES (?, ?, ?, ?)", session.ID, session.UserID, session.ExpiresAt, session.CSRFToken)
	return err
}

func (s *DBStore) Get(ctx context.Context, id string) (Session, error) {
	var session Session
	err := s.db.QueryRowContext(ctx, "SELECT id, user_id, expires_at, csrf_token FROM sessions WHERE id = ?", id).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CSRFToken)
	return session, err
}

func (s *DBStore) Save(ctx context.Context, session Session) error {
	_, err := s.db.ExecContext(ctx, "UPDATE sessions SET user_id = ?, expires_at = ?, csrf_token = ? WHERE id = ?", session.UserID, session.ExpiresAt, session.CSRFToken, session.ID)
	return err
}

func (s *DBStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", id)
	return err
}
