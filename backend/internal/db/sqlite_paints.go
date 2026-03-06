package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqlitePaintRepo struct{ db *sql.DB }

func NewPaintRepository(db *sql.DB) PaintRepository {
	return &sqlitePaintRepo{db: db}
}

func (r *sqlitePaintRepo) Create(p *models.Paint) error {
	_, err := r.db.Exec(
		`INSERT INTO paints (id, user_id, brand, name, color, paint_type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.UserID, p.Brand, p.Name, p.Color, p.PaintType, p.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create paint: %w", err)
	}
	return nil
}

func (r *sqlitePaintRepo) GetByID(id string) (*models.Paint, error) {
	p := &models.Paint{}
	err := r.db.QueryRow(
		`SELECT id, user_id, brand, name, color, paint_type, created_at FROM paints WHERE id = ?`, id,
	).Scan(&p.ID, &p.UserID, &p.Brand, &p.Name, &p.Color, &p.PaintType, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get paint: %w", err)
	}
	return p, nil
}

func (r *sqlitePaintRepo) ListByUser(userID string) ([]models.Paint, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, brand, name, color, paint_type, created_at FROM paints WHERE user_id = ? ORDER BY brand, name`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list paints: %w", err)
	}
	defer rows.Close()

	var out []models.Paint
	for rows.Next() {
		var p models.Paint
		if err := rows.Scan(&p.ID, &p.UserID, &p.Brand, &p.Name, &p.Color, &p.PaintType, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan paint: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *sqlitePaintRepo) Update(p *models.Paint) error {
	_, err := r.db.Exec(
		`UPDATE paints SET brand = ?, name = ?, color = ?, paint_type = ? WHERE id = ? AND user_id = ?`,
		p.Brand, p.Name, p.Color, p.PaintType, p.ID, p.UserID,
	)
	if err != nil {
		return fmt.Errorf("update paint: %w", err)
	}
	return nil
}

func (r *sqlitePaintRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM paints WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete paint: %w", err)
	}
	return nil
}

type sqliteMiniaturePaintRepo struct{ db *sql.DB }

func NewMiniaturePaintRepository(db *sql.DB) MiniaturePaintRepository {
	return &sqliteMiniaturePaintRepo{db: db}
}

func (r *sqliteMiniaturePaintRepo) Add(mp *models.MiniaturePaint) error {
	_, err := r.db.Exec(
		`INSERT INTO miniature_paints (id, miniature_id, paint_id, purpose, notes, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		mp.ID, mp.MiniatureID, mp.PaintID, mp.Purpose, mp.Notes, mp.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("add miniature paint: %w", err)
	}
	return nil
}

func (r *sqliteMiniaturePaintRepo) ListByMiniature(miniatureID string) ([]models.MiniaturePaint, error) {
	rows, err := r.db.Query(
		`SELECT mp.id, mp.miniature_id, mp.paint_id, mp.purpose, mp.notes, mp.created_at,
		        p.id, p.user_id, p.brand, p.name, p.color, p.paint_type, p.created_at
		 FROM miniature_paints mp
		 JOIN paints p ON p.id = mp.paint_id
		 WHERE mp.miniature_id = ?
		 ORDER BY mp.created_at ASC`, miniatureID,
	)
	if err != nil {
		return nil, fmt.Errorf("list miniature paints: %w", err)
	}
	defer rows.Close()

	var out []models.MiniaturePaint
	for rows.Next() {
		var mp models.MiniaturePaint
		var p models.Paint
		if err := rows.Scan(
			&mp.ID, &mp.MiniatureID, &mp.PaintID, &mp.Purpose, &mp.Notes, &mp.CreatedAt,
			&p.ID, &p.UserID, &p.Brand, &p.Name, &p.Color, &p.PaintType, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan miniature paint: %w", err)
		}
		mp.Paint = &p
		out = append(out, mp)
	}
	return out, rows.Err()
}

func (r *sqliteMiniaturePaintRepo) Remove(id, userID string) error {
	// Join through miniature to enforce ownership
	_, err := r.db.Exec(
		`DELETE FROM miniature_paints WHERE id = ?
		 AND miniature_id IN (SELECT id FROM miniatures WHERE user_id = ?)`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("remove miniature paint: %w", err)
	}
	return nil
}
