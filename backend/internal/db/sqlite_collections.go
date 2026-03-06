package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteCollectionRepo struct{ db *sql.DB }

func NewCollectionRepository(db *sql.DB) CollectionRepository {
	return &sqliteCollectionRepo{db: db}
}

func (r *sqliteCollectionRepo) Create(c *models.Collection) error {
	_, err := r.db.Exec(
		`INSERT INTO collections (id, user_id, name, notes, created_at) VALUES (?, ?, ?, ?, ?)`,
		c.ID, c.UserID, c.Name, c.Notes, c.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}
	return nil
}

func (r *sqliteCollectionRepo) GetByID(id string) (*models.Collection, error) {
	c := &models.Collection{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, notes, created_at FROM collections WHERE id = ?`, id,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Notes, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	return c, nil
}

func (r *sqliteCollectionRepo) ListByUser(userID string) ([]models.Collection, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, notes, created_at FROM collections WHERE user_id = ? ORDER BY name ASC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	defer rows.Close()

	var out []models.Collection
	for rows.Next() {
		var c models.Collection
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Notes, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan collection: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *sqliteCollectionRepo) Update(c *models.Collection) error {
	_, err := r.db.Exec(
		`UPDATE collections SET name = ?, notes = ? WHERE id = ? AND user_id = ?`,
		c.Name, c.Notes, c.ID, c.UserID,
	)
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}
	return nil
}

func (r *sqliteCollectionRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM collections WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}
