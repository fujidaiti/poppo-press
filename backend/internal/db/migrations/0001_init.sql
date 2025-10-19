-- moved from backend/migrations
-- schema version table
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER NOT NULL
);

-- user
CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- device
CREATE TABLE IF NOT EXISTS device (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  token_hash TEXT UNIQUE,
  last_seen_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- source
CREATE TABLE IF NOT EXISTS source (
  id INTEGER PRIMARY KEY,
  url TEXT NOT NULL UNIQUE,
  title TEXT,
  etag TEXT,
  last_modified TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT
);

-- article
CREATE TABLE IF NOT EXISTS article (
  id INTEGER PRIMARY KEY,
  source_id INTEGER NOT NULL,
  canonical_url TEXT NOT NULL,
  title TEXT NOT NULL,
  summary TEXT,
  content TEXT,
  author TEXT,
  published_at TEXT NOT NULL,
  updated_at TEXT,
  canonical_id TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (source_id) REFERENCES source(id) ON DELETE CASCADE
);

-- edition
CREATE TABLE IF NOT EXISTS edition (
  id INTEGER PRIMARY KEY,
  local_date TEXT NOT NULL UNIQUE,
  published_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- edition_article
CREATE TABLE IF NOT EXISTS edition_article (
  edition_id INTEGER NOT NULL,
  article_id INTEGER NOT NULL,
  position INTEGER NOT NULL,
  PRIMARY KEY (edition_id, article_id),
  FOREIGN KEY (edition_id) REFERENCES edition(id) ON DELETE CASCADE,
  FOREIGN KEY (article_id) REFERENCES article(id) ON DELETE CASCADE
);

-- read_state
CREATE TABLE IF NOT EXISTS read_state (
  article_id INTEGER NOT NULL,
  device_id INTEGER NOT NULL,
  is_read INTEGER NOT NULL,
  updated_at TEXT NOT NULL,
  PRIMARY KEY (article_id, device_id),
  FOREIGN KEY (article_id) REFERENCES article(id) ON DELETE CASCADE,
  FOREIGN KEY (device_id) REFERENCES device(id) ON DELETE CASCADE
);

-- bookmark
CREATE TABLE IF NOT EXISTS bookmark (
  article_id INTEGER PRIMARY KEY,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (article_id) REFERENCES article(id) ON DELETE CASCADE
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_source_url ON source(url);
CREATE INDEX IF NOT EXISTS idx_article_source_published ON article(source_id, published_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_article_canonical_id ON article(canonical_id) WHERE canonical_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_edition_local_date ON edition(local_date);
CREATE INDEX IF NOT EXISTS idx_edition_article_position ON edition_article(edition_id, position);
CREATE INDEX IF NOT EXISTS idx_read_state_device_updated ON read_state(device_id, updated_at DESC);


