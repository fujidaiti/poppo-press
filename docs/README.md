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
7. Backend development and deployment — see `../backend/README.md`
8. [Security](security.md) — auth, device tokens, storage, transport

### API Viewer

- Redoc: open `redoc.html` in a browser to view the API docs rendered from `openapi.yaml`.

### At-a-glance scope

- Single-user, self-hosted server with CLI clients
- SQLite storage (WAL), hourly feed polling; daily edition assembly
- Authentication with device-scoped tokens; multi-device support

### Instructions

- View docs:
  - Start here (this file), then follow the suggested order above.
  - API docs: open `redoc.html` or serve locally and open `http://localhost:8000/docs/redoc.html`.
    - Example: `python3 -m http.server 8000` from repo root.
  - Backend quickstart and deployment: `../backend/README.md`.
- Update API:
  - Edit `openapi.yaml`, then refresh Redoc.
- Roadmap:
  - See `../notes/ROADMAP.md` for backend-first milestones and DoD.
