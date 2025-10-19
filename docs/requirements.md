# Requirements

- Single user (self-hosted)
- Support multiple devices for the same user
- Authentication required for API and CLI access
- SQLite as the primary data store

## Tech Stack

- Frontend: CLI (Go, e.g. Cobra)
- Backend: Golang (HTTP API, e.g. chi)
- Database: SQLite (WAL enabled)

## Non-Functional Requirements

- Reliability: daily newspaper is generated at the configured local time with retries/backoff on failures.
- Fetching: poll all sources hourly (configurable) using conditional GET to distribute risk of upstream downtime.
- Performance: handle 100 registered sources by default; assemble an edition locally in under 1s.
- Storage: SQLite with WAL; backup-friendly; periodic VACUUM recommended.
- Security: password hashing (argon2id or bcrypt); device-scoped tokens; TLS recommended behind a reverse proxy.
- Portability: single static binary; supports macOS and Linux.
- Timezone: user-configurable local timezone; correct DST handling.
- Observability: structured logs; minimal metrics optional.
- Privacy: all data stored locally on the host.

## Assumptions

- Exactly one user account exists.
- System clock and timezone are correctly configured on the host.
- Sources expose standard RSS/Atom with usable `published`/`updated` times.
- Network access is available to fetch feeds.
- One edition per day; includes items published or updated in the last 24 hours at publish time.

## Non-Goals

- Multi-user or multi-tenant operation.
- Web UI.
- Push/real-time notifications.
- Cross-host sync or federation.
