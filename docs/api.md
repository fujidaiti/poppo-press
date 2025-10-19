# API

Base URL: `/v1`

## Authentication

- Login
  - POST `/auth/login`
  - Body: `{ "username": string, "password": string, "deviceName": string }`
  - 200: `{ "token": string, "deviceId": string }`
  - 401 on invalid credentials; 429 on rate limit.
- Auth Header: `Authorization: Bearer <token>`
- Logout (revoke device)
  - POST `/auth/logout` → revokes current token

## Sources

- GET `/sources` → `[ { id, url, title, createdAt } ]`
- POST `/sources` Body: `{ url: string }` → `201 { id }`
- DELETE `/sources/{id}` → `204`

## Scheduler

- Hourly job fetches all sources (conditional GET) and persists new/updated articles.

## Editions

- GET `/editions` Query: `page, pageSize` → paginated list `[ { id, localDate, publishedAt, articleCount } ]`
- GET `/editions/{id}` → `{ id, localDate, publishedAt, articles: [ ... ] }`

## Articles

- GET `/articles` Query: `editionId?, sourceId?, readState?, q?`
- GET `/articles/{id}` → article detail
- POST `/articles/{id}/read` Body: `{ isRead: boolean }` → `204`

## Read Later

- GET `/read-later` → `[ article ]`
- POST `/read-later/{id}` → `204`
- DELETE `/read-later/{id}` → `204`

## Devices

- GET `/devices` → `[ { id, name, lastSeenAt, createdAt } ]`
- DELETE `/devices/{id}` → `204` (revoke by id)

## Errors

- Error shape: `{ error: { code: string, message: string, details?: any } }`
- Common codes: `unauthorized`, `forbidden`, `not_found`, `validation_failed`, `rate_limited`, `conflict`.

## Pagination

- Query: `page` (1-based), `pageSize` (default 20, max 100)
- Response headers: `X-Total-Count`
