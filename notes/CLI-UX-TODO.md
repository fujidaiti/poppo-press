# CLI UX TODO

Use this checklist to track discussion and decisions before implementation.

## Authentication and first-run

- [ ] Device naming defaults
- [ ] Non-interactive login flags (e.g., --username/--password/--device)
- [ ] Token storage permissions and location
- [ ] Mask secrets in prompts and logs

## Command ergonomics

- [ ] Consistent verbs/nouns across commands
- [ ] Short/long flags and sensible defaults
- [ ] Confirmations for destructive ops (e.g., `source rm`, `device revoke`)

## Output and formatting

- [ ] Column widths and wrapping strategy
- [ ] Alignment and truncation rules
- [ ] Pager detection and override flag
- [ ] Color/theme toggle (and no-color)

## Filtering and sorting

- [ ] Global flags: --limit, --sort, --read-state, date ranges
- [ ] Stable sort order across runs

## Machine-readable modes

- [ ] `--json` / `--tsv` outputs
- [ ] Quiet mode (`-q`) for scripting
- [ ] Predictable exit codes per scenario

## Interactive vs non-interactive

- [ ] TTY detection to avoid prompts in CI
- [ ] `--yes/--no-confirm` non-interactive toggles
- [ ] Progress/spinners with non-TTY fallback

## Timezones and timestamps

- [ ] Display in configured local TZ
- [ ] Relative vs absolute time toggle

## Links and actions

- [ ] `--open` (open in browser) vs copy URL behavior
- [ ] Consistent article ID display and copy helpers

## Errors and diagnostics

- [ ] Clear, actionable error messages
- [ ] Retry hints for rate limits/network failures
- [ ] `--verbose` request/response trace toggle

## Caching / offline behavior

- [ ] Graceful offline behavior and messaging
- [ ] Cached listings and invalidation rules
- [ ] `cache clear` command (if needed)

## Profiles and configuration

- [ ] Multiple server profiles (e.g., `--profile`)
- [ ] Precedence: flag > env > file
- [ ] `config view/edit` commands

## Completions and help

- [ ] Rich `--help` with examples on each command
- [ ] Shell completion install/update commands

## Security UX

- [ ] Token revocation feedback and status
- [ ] Show active device in prompts
- [ ] Warn on duplicate device names

---

## Open Questions

- [ ] Default column width and wrapping policy
- [ ] Color by default or opt-in only?
- [ ] Where to store per-user cache (path and size limits)?

## Decisions

- [ ] (Fill in once agreed)
