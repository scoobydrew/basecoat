package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/drews/basecoat/internal/models"
)

type sqliteCatalogRepo struct{ db *sql.DB }

func NewCatalogRepository(db *sql.DB) CatalogRepository {
	return &sqliteCatalogRepo{db: db}
}

// FindGame searches by case-insensitive name. If publisher is non-empty it is
// also matched, which helps disambiguate games that share a name.
func (r *sqliteCatalogRepo) FindGame(name, publisher string) (*models.CatalogGame, error) {
	var (
		query string
		args  []any
	)
	if publisher != "" {
		query = `SELECT id, name, publisher, year, created_at FROM catalog_games
		          WHERE LOWER(name) = LOWER(?) AND LOWER(publisher) = LOWER(?) LIMIT 1`
		args = []any{name, publisher}
	} else {
		query = `SELECT id, name, publisher, year, created_at FROM catalog_games
		          WHERE LOWER(name) = LOWER(?) LIMIT 1`
		args = []any{name}
	}

	g := &models.CatalogGame{}
	err := r.db.QueryRow(query, args...).Scan(&g.ID, &g.Name, &g.Publisher, &g.Year, &g.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find catalog game: %w", err)
	}
	return g, nil
}

func (r *sqliteCatalogRepo) CreateGame(g *models.CatalogGame) error {
	_, err := r.db.Exec(
		`INSERT INTO catalog_games (id, name, publisher, year, created_at) VALUES (?, ?, ?, ?, ?)`,
		g.ID, g.Name, g.Publisher, g.Year, g.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create catalog game: %w", err)
	}
	return nil
}

// FindBox searches case-insensitively within a catalog game.
func (r *sqliteCatalogRepo) FindBox(catalogGameID, name string) (*models.CatalogBox, error) {
	// Normalise name for matching: collapse whitespace, lowercase.
	normalized := strings.ToLower(strings.Join(strings.Fields(name), " "))
	b := &models.CatalogBox{}
	err := r.db.QueryRow(
		`SELECT id, catalog_game_id, name, created_at FROM catalog_boxes
		  WHERE catalog_game_id = ? AND LOWER(TRIM(name)) = ? LIMIT 1`,
		catalogGameID, normalized,
	).Scan(&b.ID, &b.CatalogGameID, &b.Name, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find catalog box: %w", err)
	}
	return b, nil
}

func (r *sqliteCatalogRepo) CreateBox(b *models.CatalogBox) error {
	_, err := r.db.Exec(
		`INSERT INTO catalog_boxes (id, catalog_game_id, name, created_at) VALUES (?, ?, ?, ?)`,
		b.ID, b.CatalogGameID, b.Name, b.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create catalog box: %w", err)
	}
	return nil
}

func (r *sqliteCatalogRepo) ListBoxMiniatures(catalogBoxID string) ([]models.CatalogMiniature, error) {
	rows, err := r.db.Query(
		`SELECT id, catalog_box_id, name, unit_type, quantity, created_at
		   FROM catalog_miniatures WHERE catalog_box_id = ? ORDER BY name ASC`,
		catalogBoxID,
	)
	if err != nil {
		return nil, fmt.Errorf("list catalog miniatures: %w", err)
	}
	defer rows.Close()

	var out []models.CatalogMiniature
	for rows.Next() {
		var m models.CatalogMiniature
		if err := rows.Scan(&m.ID, &m.CatalogBoxID, &m.Name, &m.UnitType, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan catalog miniature: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *sqliteCatalogRepo) CreateMiniature(m *models.CatalogMiniature) error {
	_, err := r.db.Exec(
		`INSERT INTO catalog_miniatures (id, catalog_box_id, name, unit_type, quantity, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		m.ID, m.CatalogBoxID, m.Name, m.UnitType, m.Quantity, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create catalog miniature: %w", err)
	}
	return nil
}
