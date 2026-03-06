package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteImageRepo struct{ db *sql.DB }

func NewImageRepository(db *sql.DB) ImageRepository {
	return &sqliteImageRepo{db: db}
}

func (r *sqliteImageRepo) Create(img *models.MiniatureImage) error {
	_, err := r.db.Exec(
		`INSERT INTO miniature_images (id, miniature_id, user_id, stage, storage_path, caption, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		img.ID, img.MiniatureID, img.UserID, img.Stage, img.StoragePath, img.Caption, img.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	return nil
}

func (r *sqliteImageRepo) ListByMiniature(miniatureID string) ([]models.MiniatureImage, error) {
	rows, err := r.db.Query(
		`SELECT id, miniature_id, user_id, stage, storage_path, caption, created_at
		 FROM miniature_images WHERE miniature_id = ? ORDER BY created_at ASC`, miniatureID,
	)
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	defer rows.Close()

	var out []models.MiniatureImage
	for rows.Next() {
		var img models.MiniatureImage
		if err := rows.Scan(&img.ID, &img.MiniatureID, &img.UserID, &img.Stage, &img.StoragePath, &img.Caption, &img.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan image: %w", err)
		}
		out = append(out, img)
	}
	return out, rows.Err()
}

func (r *sqliteImageRepo) GetByID(id string) (*models.MiniatureImage, error) {
	img := &models.MiniatureImage{}
	err := r.db.QueryRow(
		`SELECT id, miniature_id, user_id, stage, storage_path, caption, created_at
		 FROM miniature_images WHERE id = ?`, id,
	).Scan(&img.ID, &img.MiniatureID, &img.UserID, &img.Stage, &img.StoragePath, &img.Caption, &img.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}
	return img, nil
}

func (r *sqliteImageRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM miniature_images WHERE id = ? AND user_id = ?`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}
