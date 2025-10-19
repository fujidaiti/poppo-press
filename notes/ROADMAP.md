# Roadmap

Backend-first plan to deliver the Go API server before the CLI.

## Milestones (backend first)

### M0 — Project bootstrap

- [x] Initialize Go module and baseline layout (`backend/cmd/server`, `backend/internal/...`)
– Add dependencies:
  - [x] `chi`
  - [x] SQLite driver (`modernc.org/sqlite`)
  - [x] `gofeed`
  - [x] password hashing (argon2id)
  - [x] scheduler (cron)
– Configuration:
  - [x] env (HTTP addr, DB path)
  - [x] file (timezone, publish time)
- [ ] Makefile or task runner; lint config; basic CI stub
- [x] Health endpoint; version/build info endpoint

Definition of Done: server starts; config loads; SQLite file writable.

### M1 — Database and migrations

- [x] Implement schema from `docs/data-model.md` using SQL migrations
- [x] Enable WAL and pragmatic PRAGMA settings at startup
- [x] Seed single user if none exists

Definition of Done: migration tool runs repeatably; schema version tracked.

### M2 — Auth and devices

- [x] Password hashing (argon2id preferred) and admin bootstrap
- [x] POST `/v1/auth/login` issues device-scoped token; store only token hash
- [x] Auth middleware validates `Authorization: Bearer` tokens
- [x] POST `/v1/auth/logout` revokes current token

Definition of Done: protected test route accessible only with valid token.

### M3 — Sources API

- [x] POST/GET/DELETE `/v1/sources` with URL validation
- [x] Initial probe fetch to verify feed; persist `etag` and `last_modified`

Definition of Done: sources can be added, listed, removed; invalid feeds rejected.

### M4 — Fetcher and hourly polling

- [x] Hourly scheduler job; per-source conditional GET with backoff on errors
- [x] Parse with `gofeed`; normalize fields; upsert articles; dedupe by canonical id
- [x] Respect per-source rate limits and timeouts

Definition of Done: hourly job persists new or updated items without duplicates.

### M5 — Aggregator and daily edition

- [ ] Assemble daily edition at configured local time (last 24h window)
- [ ] Idempotent re-runs; attach articles with ordered positions

Definition of Done: GET `/v1/editions` shows the day’s edition and counts.

### M6 — Articles and read state

- [ ] GET `/v1/articles` with filters; GET `/v1/articles/{id}` detail
- [ ] POST `/v1/articles/{id}/read` toggles device read; derive global read

Definition of Done: read state reflected in listings and detail.

### M7 — Read later

- [ ] GET/POST/DELETE `/v1/read-later`

Definition of Done: bookmarking works; duplicates prevented.

### M8 — Devices management

- [ ] GET `/v1/devices`; DELETE `/v1/devices/{id}` revoke

Definition of Done: revoked token immediately denied.

### M9 — Observability and safeguards

- [ ] Structured logging with request ids; job summaries
- [ ] Basic metrics hooks (optional); pprof guarded (optional)
- [ ] Rate limit login; sane HTTP timeouts; size limits

Definition of Done: logs are actionable; basic DoS mitigations in place.

### M10 — Deployment

- [ ] Systemd unit; environment; DB path; timezone; publish time
- [ ] Backup and restore procedure for SQLite (WAL)
- [ ] Optional Docker Compose

Definition of Done: server reliably runs on a host with documented runbook.

### M11 — CLI MVP (follow-up)

- [ ] `pp init` and `pp login`; `pp source add/list/rm`
- [ ] `pp paper read/list`; `pp later add/list/rm`; `pp device list/revoke`

Definition of Done: CLI can fully operate the backend for daily use.

## Sequence

M0 → M1 → M2 → M3 → M4 → M5 → M6 → M7 → M8 → M9 → M10 → M11

## Notes

- Use configured IANA timezone; store local-date keys for editions.
- Hourly fetch distributes risk of upstream downtime; daily assemble remains the user-facing boundary.
