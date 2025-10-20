# CLI UX

Binary name: `pp`

## Config

- Config file: `~/.config/poppo-press/config.yaml`
- Fields: `server`, `token`, `output: { pager, width }`

## Commands

### init

```console
pp init --server <url>
```

Initializes the CLI by creating a config file at `~/.config/poppo-press/config.yaml`.
This does not contact the server; it simply writes your chosen API base URL and default output settings.

Example:

```console
pp init --server http://localhost:8080
```

Resulting config (example):

```yaml
server: http://localhost:8080
token: ""            # filled by `pp login`
output:
  pager: auto
  width: 80
```

### login

```console
pp login
```

Prompts for username, password, and a device name. On success, stores the received token in your config file.
If a token already exists, it will be replaced. The device name helps you distinguish tokens across machines.

Example flow:

```console
Username: admin
Password: ********
Device name: macbook
Login successful. Token saved.
```

### source add

```console
pp source add <url>
```

Adds a new RSS/Atom source. The CLI validates the URL via the API’s probe (feed title, ETag/Last-Modified recorded).
Prints the created source id on success.

Example:

```console
$ pp source add https://example.com/feed.xml
Added source id=12 "Example Feed"
```

### source list

```console
pp source list
```

Lists all configured sources. Shows id, title, and URL.

Example:

```console
ID   TITLE                     URL
12   Example Feed              https://example.com/feed.xml
3    Another Source            https://news.example.org/rss
```

### source rm

```console
pp source rm <id>
```

Removes a source by id. This doesn’t delete already-fetched articles from the database, but new fetches will stop for that source.

Example:

```console
$ pp source rm 12
Removed source id=12
```

### paper read

```console
pp paper read [--date YYYY-MM-DD]
```

Opens the daily edition for the given date (defaults to today). Renders a numbered list of articles from the last 24 hours at the configured publish time.
You can select an article by number (implementation-specific) or open details in a follow-up command.

Example output:

```console
# Edition 2025-10-19
1. [Example Feed] Title A (08:12)
2. [Another Source] Title B (07:55)
```

### paper list

```console
pp paper list [--limit N]
```

Lists recent editions with their article counts. Useful to discover past days.

Example:

```console
ID   DATE         ARTICLES
17   2025-10-19   24
16   2025-10-18   21
15   2025-10-17   23
```

### later add

```console
pp later add <article-id>
```

Adds an article to your read-later list. Idempotent: adding the same article again is a no-op.

Example:

```console
$ pp later add 202
Added article 202 to read later
```

### later list

```console
pp later list
```

Lists articles in your read-later queue.

Example:

```console
ID   TITLE                              ADDED
202  Interesting Post                   2025-10-19 09:01
187  Deep Dive on Feeds                 2025-10-18 21:44
```

### later rm

```console
pp later rm <article-id>
```

Removes an article from your read-later list by id.

Example:

```console
$ pp later rm 202
Removed article 202 from read later
```

### device list

```console
pp device list
```

Lists devices (tokens) associated with your account, including their ids and last seen times.

Example:

```console
ID   NAME        LAST SEEN
5    macbook     2025-10-19 08:30
3    desktop     2025-10-18 22:11
```

### device revoke

```console
pp device revoke <id>
```

Revokes a device token by id. Subsequent requests using that token will be denied immediately.

Example:

```console
$ pp device revoke 5
Revoked device id=5
```

## Exit Codes

- 0 success; 1 generic error; 2 validation; 3 auth; 4 network.
