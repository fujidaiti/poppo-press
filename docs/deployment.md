# Deployment

## Binary

- Single static binary `poppo-press`.
- Flags/env: `PP_DB_PATH`, `PP_HTTP_ADDR`, `PP_TZ`, `PP_PUBLISH_TIME`.

## SQLite

- Path: e.g., `/var/lib/poppo-press/poppo.db`.
- Enable WAL, set `synchronous=NORMAL`, periodic `VACUUM`.
- File permissions: 600 for DB, 700 for parent dir.

## Systemd

- Unit (snippet):
  - `ExecStart=/usr/local/bin/poppo-press serve`
  - `Environment=PP_DB_PATH=/var/lib/poppo-press/poppo.db`
  - `Environment=PP_TZ=Europe/Berlin PP_PUBLISH_TIME=10:00`
  - `Restart=on-failure`

## Logs

- Structured logs to stdout; rotate via journald or logrotate.

## Backups

- Quiesce or use `sqlite3 .backup`; copy WAL + DB.
- Verify restore on staging.

## Docker (optional)

- Provide Compose with bind-mounted DB and config.
