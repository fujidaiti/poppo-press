# Security

## Threat Model (single-user)

- Protect against stolen devices/tokens, weak passwords, DB exfiltration.

## Authentication

- Password hashing: argon2id (preferred) or bcrypt with strong params.
- Login issues device token (random 256-bit); store only token hash server-side.
- Token in `Authorization: Bearer` header; short idle timeout optional.

## Devices

- List and revoke tokens; immediate invalidation.
- Last-seen tracking; alert on unusual IPs (future).

## Transport

- Run behind TLS-terminating reverse proxy (nginx, Caddy).
- Enforce HTTPS for remote access.

## Storage

- SQLite file permissions restrictive; backups encrypted.
- Secrets via env or protected config file.

## Abuse Controls

- Rate limit login and source add.
- Per-source fetch interval caps; backoff on errors.
