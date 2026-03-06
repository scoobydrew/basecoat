-- Drop and recreate tables to reflect the new hierarchy:
-- Collection → Game → Miniature

DROP TABLE IF EXISTS miniature_images;
DROP TABLE IF EXISTS miniature_paints;
DROP TABLE IF EXISTS miniatures;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS games;

CREATE TABLE IF NOT EXISTS collections (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniatures (
    id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    unit_type TEXT NOT NULL DEFAULT '',
    quantity INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'unpainted',
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniature_paints (
    id TEXT PRIMARY KEY,
    miniature_id TEXT NOT NULL REFERENCES miniatures(id) ON DELETE CASCADE,
    paint_id TEXT NOT NULL REFERENCES paints(id) ON DELETE CASCADE,
    purpose TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniature_images (
    id TEXT PRIMARY KEY,
    miniature_id TEXT NOT NULL REFERENCES miniatures(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stage TEXT NOT NULL DEFAULT '',
    storage_path TEXT NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_games_collection_id ON games(collection_id);
CREATE INDEX IF NOT EXISTS idx_games_user_id ON games(user_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_game_id ON miniatures(game_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_user_id ON miniatures(user_id);
CREATE INDEX IF NOT EXISTS idx_miniature_paints_miniature_id ON miniature_paints(miniature_id);
CREATE INDEX IF NOT EXISTS idx_miniature_images_miniature_id ON miniature_images(miniature_id);
CREATE INDEX IF NOT EXISTS idx_collections_user_id ON collections(user_id);
