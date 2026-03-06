package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteGameRepo struct{ db *sql.DB }

func NewGameRepository(db *sql.DB) GameRepository {
	return &sqliteGameRepo{db: db}
}

func (r *sqliteGameRepo) Create(g *models.Game) error {
	var catalogGameID *string
	if g.CatalogGameID != "" {
		catalogGameID = &g.CatalogGameID
	}
	_, err := r.db.Exec(
		`INSERT INTO games (id, collection_id, user_id, name, catalog_game_id, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		g.ID, g.CollectionID, g.UserID, g.Name, catalogGameID, g.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create game: %w", err)
	}
	return nil
}

func (r *sqliteGameRepo) GetByID(id string) (*models.Game, error) {
	g := &models.Game{}
	var (
		catalogGameID sql.NullString
		publisher     sql.NullString
		year          sql.NullInt64
	)
	err := r.db.QueryRow(`
		SELECT g.id, g.collection_id, g.user_id, g.name, g.catalog_game_id,
		       COALESCE(cg.publisher, ''), cg.year, g.created_at
		FROM games g
		LEFT JOIN catalog_games cg ON cg.id = g.catalog_game_id
		WHERE g.id = ?`, id,
	).Scan(&g.ID, &g.CollectionID, &g.UserID, &g.Name, &catalogGameID,
		&publisher, &year, &g.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get game: %w", err)
	}
	g.CatalogGameID = catalogGameID.String
	g.Publisher = publisher.String
	if year.Valid {
		y := int(year.Int64)
		g.Year = &y
	}
	return g, nil
}

func (r *sqliteGameRepo) ListByCollection(collectionID, userID string) ([]models.Game, error) {
	rows, err := r.db.Query(`
		SELECT g.id, g.collection_id, g.user_id, g.name, g.catalog_game_id,
		       COALESCE(cg.publisher, ''), cg.year, g.created_at
		FROM games g
		LEFT JOIN catalog_games cg ON cg.id = g.catalog_game_id
		WHERE g.collection_id = ? AND g.user_id = ?
		ORDER BY g.name ASC`,
		collectionID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()

	var out []models.Game
	for rows.Next() {
		var (
			g             models.Game
			catalogGameID sql.NullString
			publisher     sql.NullString
			year          sql.NullInt64
		)
		if err := rows.Scan(&g.ID, &g.CollectionID, &g.UserID, &g.Name, &catalogGameID,
			&publisher, &year, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		g.CatalogGameID = catalogGameID.String
		g.Publisher = publisher.String
		if year.Valid {
			y := int(year.Int64)
			g.Year = &y
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *sqliteGameRepo) Update(g *models.Game) error {
	_, err := r.db.Exec(
		`UPDATE games SET name = ? WHERE id = ? AND user_id = ?`,
		g.Name, g.ID, g.UserID,
	)
	if err != nil {
		return fmt.Errorf("update game: %w", err)
	}
	return nil
}

func (r *sqliteGameRepo) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM games WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}
	return nil
}
