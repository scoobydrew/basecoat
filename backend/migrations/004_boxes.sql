-- Add boxes between games and miniatures.
-- Drop and recreate miniatures since they now belong to a box.

DROP TABLE IF EXISTS miniature_images;
DROP TABLE IF EXISTS miniature_paints;
DROP TABLE IF EXISTS miniatures;

CREATE TABLE IF NOT EXISTS boxes (
    id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS miniatures (
    id TEXT PRIMARY KEY,
    box_id TEXT NOT NULL REFERENCES boxes(id) ON DELETE CASCADE,
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

CREATE INDEX IF NOT EXISTS idx_boxes_game_id ON boxes(game_id);
CREATE INDEX IF NOT EXISTS idx_boxes_user_id ON boxes(user_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_box_id ON miniatures(box_id);
CREATE INDEX IF NOT EXISTS idx_miniatures_user_id ON miniatures(user_id);
CREATE INDEX IF NOT EXISTS idx_miniature_paints_miniature_id ON miniature_paints(miniature_id);
CREATE INDEX IF NOT EXISTS idx_miniature_images_miniature_id ON miniature_images(miniature_id);
