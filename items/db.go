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

func (s *DBStore) Add(ctx context.Context, item Item) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO items (name, quantity) VALUES (?, ?)", item.Name, item.Quantity)
	return err
}

func (s *DBStore) All(ctx context.Context) ([]Item, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, quantity, created_at, updated_at FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *DBStore) Get(ctx context.Context, id int) (Item, error) {
	var item Item
	err := s.db.QueryRowContext(ctx, "SELECT id, name, quantity, created_at, updated_at FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (s *DBStore) Update(ctx context.Context, id int, item Item) error {
	_, err := s.db.ExecContext(ctx, "UPDATE items SET name = ?, quantity = ? WHERE id = ?", item.Name, item.Quantity, id)
	return err
}

func (s *DBStore) Delete(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM items WHERE id = ?", id)
	return err
}
