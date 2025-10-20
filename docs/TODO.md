## Future: Secure token storage (Keychain / Credential Store)

Enhance token handling by integrating with OS credential stores instead of plain config files:

- macOS: Keychain (Security framework)
- Linux: Secret Service (libsecret) and KWallet as fallback
- Windows: Credential Manager (DPAPI)

Requirements:

- Store/retrieve token by profile/server URL key
- Headless mode: allow fallback to file or env (`PP_TOKEN`) when keyring unavailable
- Clear token on `pp logout` across storage backends
- Avoid new runtime deps for basic use; make keyring optional and auto-detected

CLI flags/env:

- `--no-keyring` to force file-based storage
- `PP_USE_KEYRING=1` to opt-in when available

Security notes:

- Never print tokens; mask in logs/errors
- Document backup/restore implications (keychain not portable)

## Future: Mask secrets in prompts and logs

Ensure sensitive values are never displayed or persisted in plaintext:

- Redact values of secret-like flags (`--password`, `--token`) in error messages and logs
- Redact environment variables `PP_PASSWORD`, `PP_TOKEN` if echoed or logged
- Redact tokens in HTTP traces when `--verbose` is enabled (e.g., `Authorization: Bearer ****`)
- Avoid printing config contents containing `token`; show masked or omit entirely
- Add centralized redaction utility used by all logging and error formatting

## Future: Device ID generation (advanced)

Explore stronger, privacy-safe defaults for device identifiers beyond user-specified names:

- Random ID (preferred baseline)
  - Generate UUIDv7 or ULID on first login; persist locally; send as device id
  - Store canonical full id; display short form (e.g., first 12 base32 chars)
  - Pros: privacy-friendly, no platform coupling; Cons: not stable across fresh reinstalls

- App-scoped machine hash
  - deviceId = base32(sha256("poppo-press" + machineId))
  - Sources: macOS IOPlatformUUID; Linux /etc/machine-id; Windows MachineGuid
  - Never transmit raw identifiers; only send the hash; store full, show short
  - Fallback if missing: hash(hostname + username + primary MAC)
  - Pros: stable across reinstalls on same machine; Cons: privacy considerations and platform variance

- Hybrid approach
  - Default to random; support `--auto-device` to derive from machine hash for stability
  - Allow `pp device rename` to change human-friendly label without changing id

Privacy/Security considerations

- Do not include PII (no raw serials/MACs); namespace/salt to app scope
- Stable but non-trackable across apps; collision handling; normalize to `[a-z0-9-]`
- Rotation policy: allow `pp device rotate` to regenerate id and revoke old token

UX notes

- Suggest default name as `<hostname>@<shortId>`; user can override
- Keep id immutable; treat name as mutable display label

## Future: Machine-readable modes

Add explicit machine-friendly output modes and scripting aids:

- Outputs
  - `--json`: structured JSON; top-level objects/arrays; stable field names
  - `--jsonl`: line-delimited JSON for streaming lists
  - `--tsv`: tab-separated with header row by default; `--no-header` to omit
- Quiet mode
  - `-q/--quiet` suppresses human messages; only data (stdout) and errors (stderr)
- Exit codes
  - Keep CLI exit codes predictable and consistent with docs: 0 success; 1 generic; 2 validation; 3 auth; 4 network
  - No human chatter on stdout in machine modes; errors to stderr only
- Stability & selection
  - Avoid breaking field removals; additive changes only; consider `schema: "v1"` in JSON
  - `--fields title,id,date,...` to project columns/keys for TSV/JSON
- Examples (to document later)
  - `pp articles list --json | jq '.[] | select(.read==false)'`
  - `pp device list --tsv --no-header | cut -f1`
