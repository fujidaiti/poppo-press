# CLI Roadmap (Go)

Target: Ship a minimal, script-friendly Go CLI (`pp`) aligned with `docs/cli-ux.md`.
Scope: Non-interactive by default, raw text output, no colors/pager, command-scoped flags.

## C0: Foundations and constraints

- [x] Lock UX constraints per `docs/cli-ux.md` (raw output, no global flags, no pager/colors)
- [x] Decide command framework: `spf13/cobra` (no shell completions)
- [x] Directory layout and module strategy: separate Go module under `cli/`

Acceptance:

- [x] `cli/` contains `go.mod` and `main` wiring with root command
- [x] `pp --help` prints grouped help; no completion commands shown

## C1: Project scaffold

- [x] Create `cli/cmd/pp/main.go` and internal packages (auth, config, httpc, commands)
- [x] Subcommand skeletons: `init`, `login`, `source`, `paper`, `later`, `device`, `config tz`, `config tz set`
- [ ] Common flag wiring (per-command only)
- [ ] Build outputs to `cli/build/` (use `go build -o build/pp`)

Acceptance:

- [ ] `pp` runs and shows subcommands; each `--help` has examples

## C2: Config and filesystem

- [x] Read/write YAML at `~/.config/poppo-press/config.yaml` (Linux/macOS) and `%APPDATA%/Poppo Press/config.yaml` (Windows)
- [x] Fields: `server`, `token`, `timezone`, `output: { pager }`
- [x] File perms: 600; dir perms: 700 (Windows: user-only ACL)
- [x] Commands: `pp init --server <url>`, `pp config tz`, `pp config tz set <IANA-TZ>`

Acceptance:

- [x] Init writes config with correct perms; tz commands get/set properly

## C3: HTTP client

- [x] Base URL from config; bearer token injection
- [x] Timeouts and sensible defaults
- [x] `--verbose` outputs redacted request/response traces
- [x] Exit codes: 0 success; 1 generic; 2 validation; 3 auth; 4 network

Acceptance:

- [x] Network/HTTP errors map to documented exit codes; `--verbose` works

## C4: Authentication

- [x] `pp login --device <name> --username <user> --password <pass>` (non-interactive)
- [x] Store token in config; replace if exists
- [x] Messages: success, validation errors, auth failures

Acceptance:

- [x] Login persists token; respects `PP_USERNAME`/`PP_PASSWORD`/`PP_TOKEN`

## C5: Sources

- [x] `pp source add <url>` → prints created id and title
- [x] `pp source list` → raw table of id, title, URL
- [x] `pp source rm <id>` → confirms removal (prompt unless `--force`)

Acceptance:

- [x] Outputs include IDs; no formatting/wrapping beyond raw text

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

---

## Contributor Guide (CLI)

### Setup

- Install Go 1.25+
- Create a new module under `cli/`: `go mod init github.com/fujidaiti/poppo-press/cli`
- Use `spf13/cobra` (no shell completions); add other deps as needed

### Build & Run

- From `cli/`: `mkdir -p build && go build -o build/pp ./cmd/pp && ./build/pp --help`
- Point to local API: `pp init --server http://localhost:8080`
- Login (non-interactive): `PP_USERNAME=admin PP_PASSWORD=admin pp login --device devbox`

### Tests

- `go test ./...`
- Golden tests for command outputs; keep outputs raw (no colors/pager)

### UX Constraints

- Non-interactive by default; destructive prompts allow `--force`
- Command-scoped flags only; no global flags
- Raw text output; no alignment/wrapping/truncation; no colors; no pager
- Include IDs; display URLs only; never auto-open a browser
- Errors must be clear/actionable; `--verbose` prints redacted request/response traces
- Timezones: server times UTC; CLI displays in configured IANA TZ

### Code Structure

- Commands under `cli/internal/commands/...`
- Shared packages: `config`, `httpc`, `auth`
- Keep functions small and explicit; avoid hidden global state

### Branching & PRs

- One milestone or feature per PR; keep diffs small
- Update `docs/cli-ux.md` if behavior or flags change
- Update this roadmap as items complete

### Running against backend

- Start backend per `backend/README.md`
- Use a test DB; don’t commit local data files

### Security

- Never log tokens/passwords; use the redaction helper
- On revoke, device stays visible; token becomes invalid immediately
