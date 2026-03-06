package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/drews/basecoat/internal/api"
	"github.com/drews/basecoat/internal/claude"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbPath := getenv("DATABASE_PATH", "basecoat.db")
	port := getenv("PORT", "8080")
	jwtSecret := mustenv("JWT_SECRET")
	storagePath := getenv("STORAGE_PATH", "uploads")
	baseURL := getenv("BASE_URL", fmt.Sprintf("http://localhost:%s", port))
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	migrationsDir := getenv("MIGRATIONS_DIR", migrationsDirRelative())
	if err := db.Migrate(database, migrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	store, err := storage.NewLocalStorage(storagePath, baseURL+"/uploads")
	if err != nil {
		log.Fatalf("init storage: %v", err)
	}

	var claudeClient *claude.Client
	if anthropicKey != "" {
		claudeClient = claude.NewClient(anthropicKey)
	} else {
		log.Println("ANTHROPIC_API_KEY not set — Claude mini lookup disabled")
	}

	router := api.NewRouter(api.Config{
		JWTSecret:   jwtSecret,
		Users:       db.NewUserRepository(database),
		Collections: db.NewCollectionRepository(database),
		Miniatures:  db.NewMiniatureRepository(database),
		Paints:      db.NewPaintRepository(database),
		MiniPaints:  db.NewMiniaturePaintRepository(database),
		Images:      db.NewImageRepository(database),
		Storage:     store,
		Claude:      claudeClient,
		StoragePath: storagePath,
		BaseURL:     baseURL,
	})

	addr := ":" + port
	log.Printf("basecoat listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

// migrationsDirRelative resolves the migrations directory relative to the
// binary's location, so it works whether run via `go run` or as a built binary.
func migrationsDirRelative() string {
	exe, err := os.Executable()
	if err != nil {
		return "migrations"
	}
	return filepath.Join(filepath.Dir(exe), "../../migrations")
}
