# CLI UX TODO

Use this checklist to track discussion and decisions before implementation.

## Authentication and first-run

- [x] Device naming defaults
- [x] Non-interactive login flags (e.g., --username/--password/--device)
- [x] Token storage permissions and location
  - Decision: store token in config file with secure perms; support `PP_TOKEN` override; add `logout` command. Future: optional OS keychain (see docs/TODO.md).
- [ ] Mask secrets in prompts and logs

## Command ergonomics

- [x] Consistent verbs/nouns across commands
- [x] Short/long flags and sensible defaults
- [x] Confirmations for destructive ops (e.g., `source rm`, `device revoke`)

## Output and formatting

- [x] Column widths and wrapping strategy
- [x] Alignment and truncation rules
- [x] Pager detection and override flag
- [x] Color/theme toggle (and no-color)

## Filtering and sorting

- [x] Global flags: --limit, --sort, --read-state, date ranges
- [x] Stable sort order across runs

## Machine-readable modes

- [x] `--json` / `--tsv` outputs
- [x] Quiet mode (`-q`) for scripting
- [x] Predictable exit codes per scenario
  - Decision: not in v1; prefer raw text for now. Future tracked in docs/TODO.md

## Interactive vs non-interactive

- [x] TTY detection to avoid prompts in CI
- [x] `--yes/--no-confirm` non-interactive toggles
- [x] Progress/spinners with non-TTY fallback

## Timezones and timestamps

- [x] Display in configured local TZ
- [x] Relative vs absolute time toggle

## Links and actions

- [x] `--open` (open in browser) vs copy URL behavior
- [x] Consistent article ID display and copy helpers

## Errors and diagnostics

- [x] Clear, actionable error messages
- [x] Retry hints for rate limits/network failures
- [x] `--verbose` request/response trace toggle

## Caching / offline behavior

- [x] Graceful offline behavior and messaging
- [x] Cached listings and invalidation rules
- [x] `cache clear` command (if needed)

## Profiles and configuration

- [x] Multiple server profiles (e.g., `--profile`)
- [x] Precedence: flag > env > file
- [x] `config view/edit` commands

## Completions and help

- [x] Rich `--help` with examples on each command
- [x] Shell completion install/update commands

## Security UX

- [x] Token revocation feedback and status
- [x] Show active device in prompts
- [x] Warn on duplicate device names

---

## Open Questions

- [x] Default column width and wrapping policy — none; emit raw text
- [x] Color by default or opt-in only? — none; no colorization
- [ ] Where to store per-user cache (path and size limits)?

## Decisions

- [x] Use mandatory user-specified `--device <name>`; no auto-detection (early stage)
- [x] Make `login` fully non-interactive: require `--username`, `--password`, `--device`
- [x] Ergonomics: verbs→nouns, consistent global flags, confirmations with `--yes` override
- [x] Output is raw: no wrapping/truncation/width formatting to favor post-processing
- [x] Disable pager and colors by default; prefer plain text for piping/processing
- [x] No global flags; filtering/sorting delegated to shell tools and JSON output
- [x] Machine-readable modes: not in v1; plan and guarantees in docs/TODO.md
- [x] All commands non-interactive by default; destructive prompts allowed, bypass with `--force`; no spinners
- [x] DB stores timestamps in UTC (TEXT ISO8601); CLI displays in configured IANA TZ; add `pp config tz` and `pp config tz set <IANA-TZ>`
- [x] Relative time toggle: not implemented in v1; prefer absolute timestamps; future option can be added
- [x] Links: display URLs only; never auto-open a browser
- [x] IDs: include article/edition IDs in outputs for ease of reuse in commands
- [x] Errors: clear, actionable; retry hints for rate limits/network; `--verbose` flag for HTTP request/response traces
- [x] No offline mode or caching; always fetch live data
- [x] Profiles: single profile only; no `--profile`
- [x] Config management: no `config view/edit`; edit file manually
- [x] Precedence: no global precedence matrix; read single config; flags override per command; env overrides limited (e.g., `PP_TOKEN` only)
- [x] Help: rich per-command `--help` with examples; no shell completions
- [x] Revocation UX: keep device record; on revoke, clear token secret (hash) and set `revoked_at`; `device list` shows status; revoking twice returns "already revoked" (idempotent)
- [x] Prompts include active device label/id when shown; duplicate device names rejected with clear error
- [x] Pagination: per-command `--limit` and `--offset` for list commands; no global flags
- [ ] (Fill in once agreed)
