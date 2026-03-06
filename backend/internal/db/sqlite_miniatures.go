package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteMiniatureRepo struct{ db *sql.DB }

func NewMiniatureRepository(db *sql.DB) MiniatureRepository {
	return &sqliteMiniatureRepo{db: db}
}

func (r *sqliteMiniatureRepo) Create(m *models.Miniature) error {
	_, err := r.db.Exec(
		`INSERT INTO miniatures (id, box_id, user_id, name, unit_type, quantity, status, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID, m.BoxID, m.UserID, m.Name, m.UnitType, m.Quantity, m.Status, m.Notes, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create miniature: %w", err)
	}
	return nil
}

func (r *sqliteMiniatureRepo) GetByID(id string) (*models.Miniature, error) {
	m := &models.Miniature{}
	err := r.db.QueryRow(
		`SELECT id, box_id, user_id, name, unit_type, quantity, status, notes, created_at, updated_at
		 FROM miniatures WHERE id = ?`, id,
	).Scan(&m.ID, &m.BoxID, &m.UserID, &m.Name, &m.UnitType, &m.Quantity, &m.Status, &m.Notes, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get miniature: %w", err)
	}
	return m, nil
}

func (r *sqliteMiniatureRepo) ListByBox(boxID, userID string) ([]models.Miniature, error) {
	rows, err := r.db.Query(
		`SELECT id, box_id, user_id, name, unit_type, quantity, status, notes, created_at, updated_at
		 FROM miniatures WHERE box_id = ? AND user_id = ? ORDER BY name ASC`,
		boxID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list miniatures: %w", err)
	}
	defer rows.Close()

	var out []models.Miniature
	for rows.Next() {
		var m models.Miniature
		if err := rows.Scan(&m.ID, &m.BoxID, &m.UserID, &m.Name, &m.UnitType, &m.Quantity, &m.Status, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan miniature: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *sqliteMiniatureRepo) Update(m *models.Miniature) error {
	_, err := r.db.Exec(
		`UPDATE miniatures SET name = ?, unit_type = ?, quantity = ?, status = ?, notes = ?, updated_at = ?
		 WHERE id = ? AND user_id = ?`,
		m.Name, m.UnitType, m.Quantity, m.Status, m.Notes, m.UpdatedAt, m.ID, m.UserID,
	)
	if err != nil {
		return fmt.Errorf("update miniature: %w", err)
	}
	return nil
}

func (r *sqliteMiniatureRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM miniatures WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete miniature: %w", err)
	}
	return nil
}

func (r *sqliteMiniatureRepo) GetDashboardStats(userID string) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
		ByStatus: make(map[string]int),
	}

	rows, err := r.db.Query(
		`SELECT status, SUM(quantity) FROM miniatures WHERE user_id = ? GROUP BY status`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("dashboard stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan stats: %w", err)
		}
		stats.ByStatus[status] = count
		stats.TotalMinis += count
		switch models.PaintingStatus(status) {
		case models.StatusFinished:
			stats.FinishedMinis += count
		case models.StatusUnpainted, models.StatusPrimed:
			stats.UnpaintedMinis += count
		default:
			stats.InProgressMinis += count
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if stats.TotalMinis > 0 {
		stats.ShamePercent = float64(stats.UnpaintedMinis) / float64(stats.TotalMinis) * 100
	}

	recentRows, err := r.db.Query(
		`SELECT id, box_id, user_id, name, unit_type, quantity, status, notes, created_at, updated_at
		 FROM miniatures WHERE user_id = ? ORDER BY updated_at DESC LIMIT 10`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("recent activity: %w", err)
	}
	defer recentRows.Close()

	for recentRows.Next() {
		var m models.Miniature
		if err := recentRows.Scan(&m.ID, &m.BoxID, &m.UserID, &m.Name, &m.UnitType, &m.Quantity, &m.Status, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan recent: %w", err)
		}
		stats.RecentActivity = append(stats.RecentActivity, m)
	}
	return stats, recentRows.Err()
}
