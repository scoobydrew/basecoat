package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteBoxRepo struct{ db *sql.DB }

func NewBoxRepository(db *sql.DB) BoxRepository {
	return &sqliteBoxRepo{db: db}
}

func (r *sqliteBoxRepo) Create(b *models.Box) error {
	var catalogBoxID *string
	if b.CatalogBoxID != "" {
		catalogBoxID = &b.CatalogBoxID
	}
	_, err := r.db.Exec(
		`INSERT INTO boxes (id, game_id, user_id, name, catalog_box_id, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		b.ID, b.GameID, b.UserID, b.Name, catalogBoxID, b.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create box: %w", err)
	}
	return nil
}

func (r *sqliteBoxRepo) GetByID(id string) (*models.Box, error) {
	b := &models.Box{}
	var catalogBoxID sql.NullString
	err := r.db.QueryRow(
		`SELECT id, game_id, user_id, name, catalog_box_id, created_at FROM boxes WHERE id = ?`, id,
	).Scan(&b.ID, &b.GameID, &b.UserID, &b.Name, &catalogBoxID, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get box: %w", err)
	}
	b.CatalogBoxID = catalogBoxID.String
	return b, nil
}

func (r *sqliteBoxRepo) ListByGame(gameID, userID string) ([]models.Box, error) {
	rows, err := r.db.Query(
		`SELECT id, game_id, user_id, name, catalog_box_id, created_at FROM boxes WHERE game_id = ? AND user_id = ? ORDER BY name ASC`,
		gameID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list boxes: %w", err)
	}
	defer rows.Close()

	var out []models.Box
	for rows.Next() {
		var b models.Box
		var catalogBoxID sql.NullString
		if err := rows.Scan(&b.ID, &b.GameID, &b.UserID, &b.Name, &catalogBoxID, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan box: %w", err)
		}
		b.CatalogBoxID = catalogBoxID.String
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *sqliteBoxRepo) Update(b *models.Box) error {
	_, err := r.db.Exec(
		`UPDATE boxes SET name = ? WHERE id = ? AND user_id = ?`,
		b.Name, b.ID, b.UserID,
	)
	if err != nil {
		return fmt.Errorf("update box: %w", err)
	}
	return nil
}

func (r *sqliteBoxRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM boxes WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete box: %w", err)
	}
	return nil
}

func (r *sqliteBoxRepo) SetCatalogBoxID(id, catalogBoxID string) error {
	_, err := r.db.Exec(`UPDATE boxes SET catalog_box_id = ? WHERE id = ?`, catalogBoxID, id)
	if err != nil {
		return fmt.Errorf("set catalog box id: %w", err)
	}
	return nil
}
