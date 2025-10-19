# CLI UX

Binary name: `pp`

## Config

- Config file: `~/.config/poppo-press/config.yaml`
- Fields: `server`, `token`, `output: { pager, width }`

## Commands

- `pp init --server <url>`
  - Writes config with server URL.

- `pp login`
  - Prompts for username/password and device name; stores token.

- `pp source add <url>`
- `pp source list`
- `pp source rm <id>`

- `pp paper read [--date YYYY-MM-DD]`
  - Opens today’s edition by default; interactive navigation.

- `pp paper list [--limit N]`

- `pp later add <article-id>`
- `pp later list`
- `pp later rm <article-id>`

- `pp device list`
- `pp device revoke <id>`

## Output Sketches

- Sources list
  - `ID  TITLE                     URL`

- Paper read
  - `# <Edition YYYY-MM-DD>`
  - `1. [Source] Title (time) [R/L]`

## Exit Codes

- 0 success; 1 generic error; 2 validation; 3 auth; 4 network.
