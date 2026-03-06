CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collections (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    notes      TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Shared catalog: game/box/miniature definitions contributed by users.
-- No user_id — these are owned by the community.

CREATE TABLE IF NOT EXISTS catalog_games (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    publisher  TEXT NOT NULL DEFAULT '',
    year       INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS catalog_boxes (
    id              TEXT PRIMARY KEY,
    catalog_game_id TEXT NOT NULL REFERENCES catalog_games(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS catalog_miniatures (
    id             TEXT PRIMARY KEY,
    catalog_box_id TEXT NOT NULL REFERENCES catalog_boxes(id) ON DELETE CASCADE,
    name           TEXT NOT NULL,
    unit_type      TEXT NOT NULL DEFAULT '',
    quantity       INTEGER NOT NULL DEFAULT 1,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User data: links to catalog entries, tracks per-user painting state.

CREATE TABLE IF NOT EXISTS games (
    id              TEXT PRIMARY KEY,
    collection_id   TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    catalog_game_id TEXT REFERENCES catalog_games(id),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS boxes (
    id             TEXT PRIMARY KEY,
    game_id        TEXT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id        TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name           TEXT NOT NULL,
    catalog_box_id TEXT REFERENCES catalog_boxes(id),
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniatures (
    id         TEXT PRIMARY KEY,
    box_id     TEXT NOT NULL REFERENCES boxes(id) ON DELETE CASCADE,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    unit_type  TEXT NOT NULL DEFAULT '',
    quantity   INTEGER NOT NULL DEFAULT 1,
    status     TEXT NOT NULL DEFAULT 'unpainted',
    notes      TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS paints (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    brand      TEXT NOT NULL,
    name       TEXT NOT NULL,
    color      TEXT NOT NULL DEFAULT '',
    paint_type TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniature_paints (
    id           TEXT PRIMARY KEY,
    miniature_id TEXT NOT NULL REFERENCES miniatures(id) ON DELETE CASCADE,
    paint_id     TEXT NOT NULL REFERENCES paints(id) ON DELETE CASCADE,
    purpose      TEXT NOT NULL DEFAULT '',
    notes        TEXT NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniature_images (
    id           TEXT PRIMARY KEY,
    miniature_id TEXT NOT NULL REFERENCES miniatures(id) ON DELETE CASCADE,
    user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stage        TEXT NOT NULL DEFAULT '',
    storage_path TEXT NOT NULL,
    caption      TEXT NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_collections_user_id        ON collections(user_id);
CREATE INDEX IF NOT EXISTS idx_games_collection_id        ON games(collection_id);
CREATE INDEX IF NOT EXISTS idx_games_user_id              ON games(user_id);
CREATE INDEX IF NOT EXISTS idx_boxes_game_id              ON boxes(game_id);
CREATE INDEX IF NOT EXISTS idx_boxes_user_id              ON boxes(user_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_box_id          ON miniatures(box_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_user_id         ON miniatures(user_id);
CREATE INDEX IF NOT EXISTS idx_paints_user_id             ON paints(user_id);
CREATE INDEX IF NOT EXISTS idx_miniature_paints_mini_id   ON miniature_paints(miniature_id);
CREATE INDEX IF NOT EXISTS idx_miniature_images_mini_id   ON miniature_images(miniature_id);
CREATE INDEX IF NOT EXISTS idx_catalog_boxes_game_id      ON catalog_boxes(catalog_game_id);
CREATE INDEX IF NOT EXISTS idx_catalog_minis_box_id       ON catalog_miniatures(catalog_box_id);
