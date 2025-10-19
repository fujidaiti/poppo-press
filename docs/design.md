# Design

## Feature Set and Criteria

- Sources
  - Add, list, remove RSS/Atom sources.
  - Validate URL reachability and feed type on add.
  - Acceptance: adding an invalid URL fails with a clear error.

- Daily Newspaper
  - Generate one edition per day at configured local time.
  - Include items published or updated in the last 24h window.
  - Acceptance: re-running generation within the same day is idempotent.

- Archive
  - Keep all past editions; list and open any edition.
  - Acceptance: editions remain accessible after day of publication.

- Read State
  - Mark articles read/unread per device and globally.
  - Acceptance: marking read syncs across devices for the single user.

- Read Later
  - Bookmark any article; list and remove bookmarks.
  - Acceptance: duplicates are prevented by per-article uniqueness.

- Devices
  - Register device tokens on login; list and revoke.
  - Acceptance: revoked tokens are immediately denied.

- Authentication
  - Single local account; password-based login.
  - Acceptance: strong password hashing; rate limiting on login.

## User Flows

1) First Run
   - Start API → create admin user → set timezone and publish time → enable WAL.
   - CLI `pp init` stores server URL; `pp login` stores token.

2) Add Source
   - CLI `pp source add <url>` → API validates, fetches once, stores ETag/Last-Modified.

3) Read Today’s Paper
   - At publish time: scheduler fetches all sources with conditional GET; normalize items; dedupe; assemble edition; persist.
   - CLI `pp paper read` lists edition articles; navigate; mark read.

4) Read Archive
   - CLI `pp paper list` → select edition → read.

5) Read Later
   - From listing, `pp later add <article-id>`; later `pp later read`.

6) Devices
   - `pp device list`; `pp device revoke <id>` to invalidate a token.

## Edge Cases

- Empty or invalid feeds: skip with warning; do not fail the whole run.
- Duplicates across sources: dedupe by canonical URL + normalized title + published time.
- Timezone or DST change: compute window using IANA TZ; store publish clock time, not UTC instant.
- Missing `published`: fallback to `updated` or fetch time.
- Large feeds: limit items per fetch; paginate if needed.
- Network failures: exponential backoff; cap retries; continue others.

## Scheduling

- Configurable publish time in local timezone.
- At publish time: full assemble job.
- Hourly fetch job: poll sources every hour using ETag/Last-Modified; persist new/updated articles; failures isolated per source.
- Idempotency: edition key = date in local TZ; re-run overwrites same edition safely.

## Configuration

- Sources saved in DB.
- Settings: timezone, publish time, network timeouts, max sources, retry policy.
- Precedence: env > config file > defaults.
