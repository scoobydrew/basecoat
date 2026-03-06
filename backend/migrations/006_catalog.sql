CREATE TABLE catalog_games (
  id        TEXT PRIMARY KEY,
  name      TEXT NOT NULL,
  publisher TEXT NOT NULL DEFAULT '',
  year      INTEGER,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalog_boxes (
  id              TEXT PRIMARY KEY,
  catalog_game_id TEXT NOT NULL REFERENCES catalog_games(id) ON DELETE CASCADE,
  name            TEXT NOT NULL,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalog_miniatures (
  id              TEXT PRIMARY KEY,
  catalog_box_id  TEXT NOT NULL REFERENCES catalog_boxes(id) ON DELETE CASCADE,
  name            TEXT NOT NULL,
  unit_type       TEXT NOT NULL DEFAULT '',
  quantity        INTEGER NOT NULL DEFAULT 1,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE games ADD COLUMN catalog_game_id TEXT REFERENCES catalog_games(id);
ALTER TABLE boxes ADD COLUMN catalog_box_id  TEXT REFERENCES catalog_boxes(id);
