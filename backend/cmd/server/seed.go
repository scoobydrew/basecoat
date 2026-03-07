package main

import (
	"log"
	"time"

	"github.com/drews/basecoat/internal/auth"
	"github.com/drews/basecoat/internal/db"
	"github.com/drews/basecoat/internal/models"
	"github.com/google/uuid"
)

// seedDevData populates the database with test data if no users exist yet.
// Safe to call on every startup — it's a no-op after the first run.
func seedDevData(repos db.Repos) {
	users, err := repos.Users.ListAll()
	if err != nil || len(users) > 0 {
		return
	}

	log.Println("seeding dev data…")

	now := time.Now()

	// ── User ──────────────────────────────────────────────────────────────────
	hash, err := auth.HashPassword("password")
	if err != nil {
		log.Printf("seed: hash password: %v", err)
		return
	}
	user := &models.User{
		ID:           uuid.NewString(),
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: hash,
		CreatedAt:    now,
	}
	if err := repos.Users.Create(user); err != nil {
		log.Printf("seed: create user: %v", err)
		return
	}

	// ── Collection ────────────────────────────────────────────────────────────
	col := &models.Collection{
		ID: uuid.NewString(), UserID: user.ID,
		Name: "My Collection", Notes: "Dev seed data", CreatedAt: now,
	}
	if err := repos.Collections.Create(col); err != nil {
		log.Printf("seed: create collection: %v", err)
		return
	}

	// ── Game: Blood Rage ──────────────────────────────────────────────────────
	year2015 := 2015
	cgBloodRage := &models.CatalogGame{
		ID: uuid.NewString(), Name: "Blood Rage",
		Publisher: "CMON", Year: &year2015, CreatedAt: now,
	}
	if err := repos.Catalog.CreateGame(cgBloodRage); err != nil {
		log.Printf("seed: catalog game: %v", err)
		return
	}

	gBloodRage := &models.Game{
		ID: uuid.NewString(), CollectionID: col.ID, UserID: user.ID,
		Name: "Blood Rage", CatalogGameID: cgBloodRage.ID, CreatedAt: now,
	}
	if err := repos.Games.Create(gBloodRage); err != nil {
		log.Printf("seed: game: %v", err)
		return
	}

	// Box: Core Set
	cbCore := &models.CatalogBox{
		ID: uuid.NewString(), CatalogGameID: cgBloodRage.ID, Name: "Core Set", CreatedAt: now,
	}
	if err := repos.Catalog.CreateBox(cbCore); err != nil {
		log.Printf("seed: catalog box: %v", err)
		return
	}

	bCore := &models.Box{
		ID: uuid.NewString(), GameID: gBloodRage.ID, UserID: user.ID,
		Name: "Core Set", CatalogBoxID: cbCore.ID, CreatedAt: now,
	}
	if err := repos.Boxes.Create(bCore); err != nil {
		log.Printf("seed: box: %v", err)
		return
	}

	coreMinis := []struct {
		name     string
		unitType string
		qty      int
		status   models.PaintingStatus
	}{
		{"Odin", "god", 1, models.StatusFinished},
		{"Thor", "god", 1, models.StatusDetailed},
		{"Tyr", "god", 1, models.StatusShaded},
		{"Heimdall", "god", 1, models.StatusBasecoated},
		{"Freya", "god", 1, models.StatusPrimed},
		{"Viking Warriors", "infantry", 16, models.StatusBasecoated},
		{"Valkyrie", "hero", 1, models.StatusPrimed},
		{"Sea Serpent", "monster", 1, models.StatusUnpainted},
		{"Fenrir Wolf", "monster", 1, models.StatusUnpainted},
		{"Midgard Serpent", "monster", 1, models.StatusUnpainted},
	}
	for _, m := range coreMinis {
		repos.Catalog.CreateMiniature(&models.CatalogMiniature{ //nolint:errcheck
			ID: uuid.NewString(), CatalogBoxID: cbCore.ID,
			Name: m.name, UnitType: m.unitType, Quantity: m.qty, CreatedAt: now,
		})
		repos.Miniatures.Create(&models.Miniature{ //nolint:errcheck
			ID: uuid.NewString(), BoxID: bCore.ID, UserID: user.ID,
			Name: m.name, UnitType: m.unitType, Quantity: m.qty,
			Status: m.status, CreatedAt: now, UpdatedAt: now,
		})
	}

	// Box: Mystics of Midgard
	cbMystics := &models.CatalogBox{
		ID: uuid.NewString(), CatalogGameID: cgBloodRage.ID, Name: "Mystics of Midgard", CreatedAt: now,
	}
	repos.Catalog.CreateBox(cbMystics) //nolint:errcheck

	bMystics := &models.Box{
		ID: uuid.NewString(), GameID: gBloodRage.ID, UserID: user.ID,
		Name: "Mystics of Midgard", CatalogBoxID: cbMystics.ID, CreatedAt: now,
	}
	repos.Boxes.Create(bMystics) //nolint:errcheck

	mysticsMinis := []struct {
		name     string
		unitType string
		qty      int
	}{
		{"Mystic Seer", "hero", 1},
		{"Cursed Warrior", "infantry", 4},
		{"Chaos Mage", "hero", 1},
	}
	for _, m := range mysticsMinis {
		repos.Catalog.CreateMiniature(&models.CatalogMiniature{ //nolint:errcheck
			ID: uuid.NewString(), CatalogBoxID: cbMystics.ID,
			Name: m.name, UnitType: m.unitType, Quantity: m.qty, CreatedAt: now,
		})
		repos.Miniatures.Create(&models.Miniature{ //nolint:errcheck
			ID: uuid.NewString(), BoxID: bMystics.ID, UserID: user.ID,
			Name: m.name, UnitType: m.unitType, Quantity: m.qty,
			Status: models.StatusUnpainted, CreatedAt: now, UpdatedAt: now,
		})
	}

	// ── Game: Warhammer 40K ───────────────────────────────────────────────────
	cgW40k := &models.CatalogGame{
		ID: uuid.NewString(), Name: "Warhammer 40,000",
		Publisher: "Games Workshop", CreatedAt: now,
	}
	repos.Catalog.CreateGame(cgW40k) //nolint:errcheck

	gW40k := &models.Game{
		ID: uuid.NewString(), CollectionID: col.ID, UserID: user.ID,
		Name: "Warhammer 40,000", CatalogGameID: cgW40k.ID, CreatedAt: now,
	}
	repos.Games.Create(gW40k) //nolint:errcheck

	cbIntercessors := &models.CatalogBox{
		ID: uuid.NewString(), CatalogGameID: cgW40k.ID, Name: "Intercessors", CreatedAt: now,
	}
	repos.Catalog.CreateBox(cbIntercessors) //nolint:errcheck

	bIntercessors := &models.Box{
		ID: uuid.NewString(), GameID: gW40k.ID, UserID: user.ID,
		Name: "Intercessors", CatalogBoxID: cbIntercessors.ID, CreatedAt: now,
	}
	repos.Boxes.Create(bIntercessors) //nolint:errcheck

	intercessorMinis := []struct {
		name   string
		qty    int
		status models.PaintingStatus
	}{
		{"Intercessor Sergeant", 1, models.StatusBasecoated},
		{"Intercessors", 4, models.StatusPrimed},
	}
	for _, m := range intercessorMinis {
		repos.Catalog.CreateMiniature(&models.CatalogMiniature{ //nolint:errcheck
			ID: uuid.NewString(), CatalogBoxID: cbIntercessors.ID,
			Name: m.name, UnitType: "infantry", Quantity: m.qty, CreatedAt: now,
		})
		repos.Miniatures.Create(&models.Miniature{ //nolint:errcheck
			ID: uuid.NewString(), BoxID: bIntercessors.ID, UserID: user.ID,
			Name: m.name, UnitType: "infantry", Quantity: m.qty,
			Status: m.status, CreatedAt: now, UpdatedAt: now,
		})
	}

	log.Println("seed: done — login with test@test.com / password")
}
