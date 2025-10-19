# Poppo Press

A cli-based RSS reader that allows you to read newspapers from your favorite publishers.

## Documentation Overview

This directory is the entry point for product and technical documentation.

### Suggested reading order

1. [Requirements](requirements.md) — scope, goals, assumptions, non-goals
2. [Design](design.md) — features, flows, edge cases, scheduling
3. [Architecture](architecture.md) — components, data flow, idempotency
4. [Data Model](data-model.md) — tables, indexes, invariants, migrations
5. [API](api.md) — endpoints, requests/responses, errors, pagination
6. [CLI UX](cli-ux.md) — commands, flags, output, exit codes
7. [Deployment](deployment.md) — binary, config, systemd, backups, Docker
8. [Security](security.md) — auth, device tokens, storage, transport

### At-a-glance scope

- Single-user, self-hosted server with CLI clients
- SQLite storage (WAL), hourly feed polling; daily edition assembly
- Authentication with device-scoped tokens; multi-device support
