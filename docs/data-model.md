# Data Model

## Tables

- user
  - id (PK)
  - username (unique)
  - password_hash
  - created_at

- device
  - id (PK)
  - name
  - token_hash (unique)
  - last_seen_at
  - created_at

- source
  - id (PK)
  - url (unique)
  - title
  - etag
  - last_modified
  - created_at
  - updated_at

- article
  - id (PK)
  - source_id (FK → source.id)
  - canonical_url
  - title
  - summary
  - content
  - author
  - published_at
  - updated_at
  - canonical_id (unique nullable)  
    // canonical hash from url+title+published when feed lacks GUID
  - created_at

- edition
  - id (PK)
  - local_date (YYYY-MM-DD, unique)  
    // derived from configured timezone
  - published_at (timestamp)
  - created_at

- edition_article
  - edition_id (FK → edition.id, composite PK)
  - article_id (FK → article.id, composite PK)
  - position (int)

- read_state
  - article_id (FK → article.id, composite PK)
  - device_id (FK → device.id, composite PK)
  - is_read (bool)
  - updated_at

- bookmark
  - article_id (FK → article.id, PK)
  - created_at

## Indexes

- source(url)
- article(source_id, published_at DESC)
- article(canonical_id) unique where not null
- edition(local_date) unique
- edition_article(edition_id, position)
- read_state(device_id, updated_at DESC)

## Invariants

- One edition per local_date.
- An article appears at most once in an edition.
- Bookmark uniqueness by article_id.
- Read state is per device; global read derived by any device read.

## Migration Strategy

- SQL migrations applied at startup; schema version table.
- Seed single user if none exists.
- Enable SQLite WAL and set pragmatic defaults (journal_size_limit, synchronous=NORMAL).
