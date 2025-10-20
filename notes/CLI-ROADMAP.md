# CLI Roadmap (Go)

Target: Ship a minimal, script-friendly Go CLI (`pp`) aligned with `docs/cli-ux.md`.
Scope: Non-interactive by default, raw text output, no colors/pager, command-scoped flags.

## C0: Foundations and constraints

- [ ] Lock UX constraints per `docs/cli-ux.md` (raw output, no global flags, no pager/colors)
- [ ] Decide command framework (prefer `spf13/cobra`, no shell completions enabled)
- [ ] Directory layout and module strategy (separate Go module under `cli/`)

Acceptance:

- [ ] `cli/` contains `go.mod` and `main` wiring with root command
- [ ] `pp --help` prints grouped help; no completion commands shown

## C1: Project scaffold

- [ ] Create `cli/cmd/pp/main.go` and internal packages (auth, config, httpc, commands)
- [ ] Subcommand skeletons: `init`, `login`, `source`, `paper`, `later`, `device`, `config tz`, `config tz set`
- [ ] Common flag wiring (per-command only)

Acceptance:

- [ ] `pp` runs and shows subcommands; each `--help` has examples

## C2: Config and filesystem

- [ ] Read/write YAML at `~/.config/poppo-press/config.yaml` (Linux/macOS) and `%APPDATA%/Poppo Press/config.yaml` (Windows)
- [ ] Fields: `server`, `token`, `timezone`, `output: { pager }`
- [ ] File perms: 600; dir perms: 700 (Windows: user-only ACL)
- [ ] Commands: `pp init --server <url>`, `pp config tz`, `pp config tz set <IANA-TZ>`

Acceptance:

- [ ] Init writes config with correct perms; tz commands get/set properly

## C3: HTTP client

- [ ] Base URL from config; bearer token injection
- [ ] Timeouts and sensible defaults
- [ ] `--verbose` outputs redacted request/response traces
- [ ] Exit codes: 0 success; 1 generic; 2 validation; 3 auth; 4 network

Acceptance:

- [ ] Network/HTTP errors map to documented exit codes; `--verbose` works

## C4: Authentication

- [ ] `pp login --device <name> --username <user> --password <pass>` (non-interactive)
- [ ] Store token in config; replace if exists
- [ ] Messages: success, validation errors, auth failures

Acceptance:

- [ ] Login persists token; respects `PP_USERNAME`/`PP_PASSWORD`/`PP_TOKEN`

## C5: Sources

- [ ] `pp source add <url>` → prints created id and title
- [ ] `pp source list` → raw table of id, title, URL
- [ ] `pp source rm <id>` → confirms removal (prompt unless `--force`)

Acceptance:

- [ ] Outputs include IDs; no formatting/wrapping beyond raw text

## C6: Papers (editions)

- [ ] `pp paper read [--date YYYY-MM-DD]` → numbered list with article IDs
- [ ] `pp paper list [--limit N] [--offset N]` → recent editions with counts

Acceptance:

- [ ] Dates reflect configured timezone; IDs shown for reuse

## C7: Read-later

- [ ] `pp later add <article-id>` (idempotent)
- [ ] `pp later list [--limit N] [--offset N]`
- [ ] `pp later rm <article-id>`

Acceptance:

- [ ] Messages are clear and raw; exit codes consistent

## C8: Devices

- [ ] `pp device list` (raw list; includes id, name, last seen, revoked status)
- [ ] `pp device revoke <id>` (idempotent; outcome: revoked/already revoked/not found)
- [ ] Prompts include active device label; `--force` bypasses

Acceptance:

- [ ] Revoked devices stay visible; second revoke reports already revoked

## C9: Pagination and selection

- [ ] Add per-command `--limit` and `--offset` where applicable (paper list, later list)
- [ ] Keep flags command-scoped; no global pagination flags

Acceptance:

- [ ] Pagination works; defaults documented in `--help`

## C10: Errors and diagnostics polish

- [ ] Central error formatter with retry hints for rate limits/network issues
- [ ] Redaction utility for logs/traces (tokens, passwords)

Acceptance:

- [ ] Representative failures show actionable messages and correct exit codes

## C11: Testing

- [ ] Unit tests for config, http client, flag parsing
- [ ] Golden tests for command outputs (raw text)
- [ ] Integration tests against a local test server (mock/fake)

Acceptance:

- [ ] CI runs tests headlessly; deterministic goldens

## C12: Packaging (optional early, required before release)

- [ ] Cross-compiles for macOS/Linux/Windows (amd64/arm64)
- [ ] Archive layout and checksums
- [ ] Homebrew/Scoop/Tarball instructions (docs only; no shell completions)

Acceptance:

- [ ] Reproducible builds; install docs validated

## C13: Future (tracked in docs/TODO.md)

- [ ] Machine-readable modes (`--json`, `--jsonl`, `--tsv`), `-q/--quiet`, `--fields`
- [ ] Secure token storage via OS keychain (opt-in)
- [ ] Advanced device ID generation (random/ULID or app-scoped hash)
