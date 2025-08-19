package items

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

func (s *DBStore) Add(ctx context.Context, item Item, userID int) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO items (name, quantity, user_id) VALUES (?, ?, ?)", item.Name, item.Quantity, userID)
	return err
}

func (s *DBStore) All(ctx context.Context, userID int) ([]Item, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, quantity, user_id, created_at, updated_at FROM items WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.UserID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *DBStore) Get(ctx context.Context, id int, userID int) (Item, error) {
	var item Item
	err := s.db.QueryRowContext(ctx, "SELECT id, name, quantity, user_id, created_at, updated_at FROM items WHERE id = ? AND user_id = ?", id, userID).Scan(&item.ID, &item.Name, &item.Quantity, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (s *DBStore) Update(ctx context.Context, id int, item Item, userID int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE items SET name = ?, quantity = ? WHERE id = ? AND user_id = ?", item.Name, item.Quantity, id, userID)
	return err
}

func (s *DBStore) Delete(ctx context.Context, id int, userID int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM items WHERE id = ? AND user_id = ?", id, userID)
	return err
}
