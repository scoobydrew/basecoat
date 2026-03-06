package db

import "github.com/drews/basecoat/internal/models"

// UserRepository defines all user persistence operations.
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

// CollectionRepository defines collection persistence operations.
type CollectionRepository interface {
	Create(c *models.Collection) error
	GetByID(id string) (*models.Collection, error)
	ListByUser(userID string) ([]models.Collection, error)
	Update(c *models.Collection) error
	Delete(id, userID string) error
}

// GameRepository defines game persistence operations.
type GameRepository interface {
	Create(g *models.Game) error
	GetByID(id string) (*models.Game, error)
	ListByCollection(collectionID, userID string) ([]models.Game, error)
	Update(g *models.Game) error
	Delete(id, userID string) error
}

// BoxRepository defines box persistence operations.
type BoxRepository interface {
	Create(b *models.Box) error
	GetByID(id string) (*models.Box, error)
	ListByGame(gameID, userID string) ([]models.Box, error)
	Update(b *models.Box) error
	Delete(id, userID string) error
	SetCatalogBoxID(id, catalogBoxID string) error
}

// CatalogRepository manages the shared game/box/miniature catalog.
type CatalogRepository interface {
	FindGame(name, publisher string) (*models.CatalogGame, error)
	CreateGame(g *models.CatalogGame) error
	FindBox(catalogGameID, name string) (*models.CatalogBox, error)
	CreateBox(b *models.CatalogBox) error
	ListBoxMiniatures(catalogBoxID string) ([]models.CatalogMiniature, error)
	CreateMiniature(m *models.CatalogMiniature) error
}

// MiniatureRepository defines miniature persistence operations.
type MiniatureRepository interface {
	Create(m *models.Miniature) error
	GetByID(id string) (*models.Miniature, error)
	ListByBox(boxID, userID string) ([]models.Miniature, error)
	Update(m *models.Miniature) error
	Delete(id, userID string) error
	GetDashboardStats(userID string) (*models.DashboardStats, error)
}

// PaintRepository defines paint persistence operations.
type PaintRepository interface {
	Create(p *models.Paint) error
	GetByID(id string) (*models.Paint, error)
	ListByUser(userID string) ([]models.Paint, error)
	Update(p *models.Paint) error
	Delete(id, userID string) error
}

// MiniaturePaintRepository links miniatures to paints.
type MiniaturePaintRepository interface {
	Add(mp *models.MiniaturePaint) error
	ListByMiniature(miniatureID string) ([]models.MiniaturePaint, error)
	Remove(id, userID string) error
}

// ImageRepository defines image metadata persistence.
type ImageRepository interface {
	Create(img *models.MiniatureImage) error
	ListByMiniature(miniatureID string) ([]models.MiniatureImage, error)
	GetByID(id string) (*models.MiniatureImage, error)
	Delete(id, userID string) error
}
