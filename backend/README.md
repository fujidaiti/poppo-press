# Poppo Press Backend

Single-user Go API server with SQLite storage, hourly feed polling, and daily edition assembly.

## Quickstart

Run the API locally:

```bash
cd backend
export PP_ADMIN_PASS='your-admin-pass'   # required only on first run to seed admin
export PP_HTTP_ADDR=':8080'              # optional (default :8080)
export PP_DB_PATH='poppo.db'             # optional (default poppo.db)
go run ./cmd/server
```

Health checks:

```bash
curl -s http://localhost:8080/health
curl -s http://localhost:8080/version
```

Login and call a protected route:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"your-admin-pass","deviceName":"dev"}' | jq -r .token)

curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/v1/protected/ping
```

## Configuration

Environment variables:

- `PP_HTTP_ADDR` (default `:8080`)
- `PP_DB_PATH`   (default `poppo.db`)
- `PP_TZ`        (default `Local`)
- `PP_PUBLISH_TIME` (default `08:00`)
- First run only: `PP_ADMIN_PASS` (required), `PP_ADMIN_USER` (default `admin`)

## Database

- SQLite with WAL; pragmatic PRAGMAs enabled on open
- Schema migrations run at startup
- First run seeds an admin user if the `user` table is empty

DB location is `PP_DB_PATH`. Ensure the file is writable by the process.

## API Surface (high-level)

- Auth: login/logout, device-scoped tokens
- Sources: add/list/delete with initial probe (stores ETag/Last-Modified)
- Fetcher: hourly conditional GET, parse via `gofeed`, upsert articles
- Editions: daily assembly at local `PP_PUBLISH_TIME` (last 24h window)
- Articles: list/detail; per-device read toggle; filters
- Read Later: add/list/remove (idempotent add)
- Devices: list and revoke

See detailed shapes in `../docs/api.md` (or open `../docs/redoc.html`).

## Scheduler

- Hourly: fetch sources (conditional GET)
- Daily: assemble edition at `PP_PUBLISH_TIME` in `PP_TZ`

## Observability & Safeguards

- JSON request logs with request IDs to stdout
- 1MB request body limit; sane HTTP timeouts
- Per-IP rate limit on login

## Development

Run tests:

```bash
cd backend
go test ./...
```

Build binary:

```bash
cd backend
go build -o poppo-press ./cmd/server
PP_ADMIN_PASS='your-admin-pass' ./poppo-press
```

## Deployment

### Systemd

Unit (complete example):

```
[Unit]
Description=Poppo Press API
After=network.target

[Service]
User=poppo
Group=poppo
Environment=PP_DB_PATH=/var/lib/poppo-press/poppo.db
Environment=PP_HTTP_ADDR=:8080
Environment=PP_TZ=Europe/Berlin
Environment=PP_PUBLISH_TIME=08:00
ExecStart=/usr/local/bin/poppo-press
Restart=on-failure
RestartSec=2s
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

### Logs

- Structured JSON logs to stdout; rotate via journald or logrotate.

### Backups

- Quiesce or use `sqlite3 .backup`; copy WAL + DB.
- Verify restore on staging.

### Docker (optional)

Provide Compose with bind-mounted DB and config.

```
services:
  app:
    image: ghcr.io/you/poppo-press:latest
    container_name: poppo-press
    environment:
      - PP_HTTP_ADDR=:8080
      - PP_DB_PATH=/data/poppo.db
      - PP_TZ=UTC
      - PP_PUBLISH_TIME=08:00
    volumes:
      - ./data:/data
    ports:
      - "8080:8080"
    restart: unless-stopped
```
