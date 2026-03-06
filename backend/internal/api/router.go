package api

import (
	"net/http"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/claude"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	JWTSecret   string
	Users       db.UserRepository
	Games       db.GameRepository
	Collections db.CollectionRepository
	Boxes       db.BoxRepository
	Miniatures  db.MiniatureRepository
	Paints      db.PaintRepository
	MiniPaints  db.MiniaturePaintRepository
	Images      db.ImageRepository
	Catalog     db.CatalogRepository
	Storage     storage.Storage
	Claude      *claude.Client
	StoragePath string
	BaseURL     string
}

func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.StoragePath))))

	ah := &authHandler{users: cfg.Users, jwtSecret: cfg.JWTSecret}
	r.Post("/api/auth/register", ah.register)
	r.Post("/api/auth/login", ah.login)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(cfg.JWTSecret))

		dh := &dashboardHandler{miniatures: cfg.Miniatures}
		r.Get("/api/dashboard", dh.get)

		ch := &collectionHandler{collections: cfg.Collections}
		r.Get("/api/collections", ch.list)
		r.Post("/api/collections", ch.create)
		r.Get("/api/collections/{id}", ch.get)
		r.Put("/api/collections/{id}", ch.update)
		r.Delete("/api/collections/{id}", ch.delete)

		gh := &gameHandler{games: cfg.Games, collections: cfg.Collections, catalog: cfg.Catalog}
		r.Get("/api/collections/{collectionID}/games", gh.list)
		r.Post("/api/collections/{collectionID}/games", gh.create)
		r.Get("/api/games/{id}", gh.get)
		r.Put("/api/games/{id}", gh.update)
		r.Delete("/api/games/{id}", gh.delete)

		bh := &boxHandler{boxes: cfg.Boxes, games: cfg.Games, miniatures: cfg.Miniatures, catalog: cfg.Catalog, claude: cfg.Claude}
		r.Get("/api/games/{gameID}/boxes", bh.list)
		r.Post("/api/games/{gameID}/boxes", bh.create)
		r.Post("/api/boxes/{id}/confirm", bh.confirm)
		r.Get("/api/boxes/{id}", bh.get)
		r.Put("/api/boxes/{id}", bh.update)
		r.Delete("/api/boxes/{id}", bh.delete)

		mh := &miniatureHandler{miniatures: cfg.Miniatures, paints: cfg.MiniPaints, images: cfg.Images}
		r.Get("/api/boxes/{boxID}/miniatures", mh.list)
		r.Post("/api/boxes/{boxID}/miniatures", mh.create)
		r.Get("/api/miniatures/{id}", mh.get)
		r.Patch("/api/miniatures/{id}", mh.update)
		r.Delete("/api/miniatures/{id}", mh.delete)
		r.Post("/api/miniatures/{id}/paints", mh.addPaint)
		r.Delete("/api/miniatures/{id}/paints/{paintLinkID}", mh.removePaint)
		r.Post("/api/miniatures/{id}/images", ih(cfg).upload)
		r.Delete("/api/miniatures/{id}/images/{imageID}", ih(cfg).delete)

		ph := &paintHandler{paints: cfg.Paints}
		r.Get("/api/paints", ph.list)
		r.Post("/api/paints", ph.create)
		r.Put("/api/paints/{id}", ph.update)
		r.Delete("/api/paints/{id}", ph.delete)
	})

	return r
}

func ih(cfg Config) *imageHandler {
	return &imageHandler{images: cfg.Images, miniatures: cfg.Miniatures, storage: cfg.Storage}
}
