package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Collection struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type Game struct {
	ID            string    `json:"id"`
	CollectionID  string    `json:"collection_id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	CatalogGameID string    `json:"catalog_game_id,omitempty"`
	// Populated via JOIN with catalog_games — not stored on this table.
	Publisher string `json:"publisher,omitempty"`
	Year      *int   `json:"year,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PaintingStatus string

const (
	StatusUnpainted   PaintingStatus = "unpainted"
	StatusPrimed      PaintingStatus = "primed"
	StatusBasecoated  PaintingStatus = "basecoated"
	StatusShaded      PaintingStatus = "shaded"
	StatusDetailed    PaintingStatus = "detailed"
	StatusFinished    PaintingStatus = "finished"
)

type Box struct {
	ID            string    `json:"id"`
	GameID        string    `json:"game_id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	CatalogBoxID  string    `json:"catalog_box_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Catalog types — shared across all users, no user_id.

type CatalogGame struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Publisher string    `json:"publisher"`
	Year      *int      `json:"year,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CatalogBox struct {
	ID            string    `json:"id"`
	CatalogGameID string    `json:"catalog_game_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
}

type CatalogMiniature struct {
	ID           string    `json:"id"`
	CatalogBoxID string    `json:"catalog_box_id"`
	Name         string    `json:"name"`
	UnitType     string    `json:"unit_type"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
}

type Miniature struct {
	ID    string `json:"id"`
	BoxID string `json:"box_id"`
	UserID string `json:"user_id"`
	Name         string         `json:"name"`
	UnitType     string         `json:"unit_type"`
	Quantity     int            `json:"quantity"`
	Status       PaintingStatus `json:"status"`
	Notes        string         `json:"notes"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Populated via JOIN, not stored directly
	Paints []MiniaturePaint `json:"paints,omitempty"`
	Images []MiniatureImage `json:"images,omitempty"`
}

type Paint struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Brand     string    `json:"brand"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	PaintType string    `json:"type"` // base, shade, highlight, contrast, technical, etc.
	CreatedAt time.Time `json:"created_at"`
}

type MiniaturePaint struct {
	ID          string    `json:"id"`
	MiniatureID string    `json:"miniature_id"`
	PaintID     string    `json:"paint_id"`
	Purpose     string    `json:"purpose"` // e.g. "base coat", "shade", "edge highlight"
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`

	// Populated via JOIN
	Paint *Paint `json:"paint,omitempty"`
}

type MiniatureImage struct {
	ID          string         `json:"id"`
	MiniatureID string         `json:"miniature_id"`
	UserID      string         `json:"user_id"`
	Stage       PaintingStatus `json:"stage"`
	StoragePath string         `json:"-"` // internal path, not exposed
	URL         string         `json:"url"`
	Caption     string         `json:"caption"`
	CreatedAt   time.Time      `json:"created_at"`
}

// DashboardStats is returned by the dashboard endpoint.
type DashboardStats struct {
	TotalMinis      int            `json:"total_minis"`
	FinishedMinis   int            `json:"finished_minis"`
	InProgressMinis int            `json:"in_progress_minis"`
	UnpaintedMinis  int            `json:"unpainted_minis"`
	ShamePercent    float64        `json:"shame_percent"` // % that are unpainted or primed
	ByStatus        map[string]int `json:"by_status"`
	RecentActivity  []Miniature    `json:"recent_activity"`
}
