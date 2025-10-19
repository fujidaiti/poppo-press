# Architecture

## Components

- API Server (Go/HTTP)
  - Exposes REST endpoints for sources, editions, articles, read-later, devices, auth.
  - Auth middleware validates device tokens.

- Scheduler
  - Triggers hourly fetch and daily edition assembly.
  - Uses cron-like scheduler; jobs are idempotent.

- Fetcher
  - Performs conditional GETs using ETag/Last-Modified.
  - Parses RSS/Atom, normalizes fields.

- Aggregator
  - Dedupe items across sources.
  - Builds the day’s edition and persists relationships.

- Storage (SQLite)
  - SQL migrations; WAL; indices for lookups.

- CLI
  - Interacts with API; renders lists and article content.

## Data Flow

1. Hourly: scheduler triggers fetch job; for each source: conditional GET → parse → normalize → upsert articles.
2. Daily: scheduler assembles edition from last 24h articles; dedupe and persist relationships.
3. Expose via API; CLI consumes.

## Concurrency & Idempotency

- Single-process lock on assemble job to avoid overlap.
- Edition key derived from local date; re-runs replace same edition atomically.
- Upserts keyed by `canonical_id` to avoid duplicates.

## Caching & HTTP Semantics

- Respect ETag/Last-Modified; backoff on `429`/`5xx`.
- Per-source fetch interval caps.

## Configuration

- Via env and config file: timezone, publish time, DB path, network timeouts.

## Error Handling

- Per-source isolation: failures do not abort entire assemble job.
- Structured logs with context; summaries after job completes.
