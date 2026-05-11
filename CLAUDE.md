# CLAUDE.md — lgtv-cli

Context for future Claude sessions working in this repo.

## What this is

A thin Go client for [`lgtv-core`](https://github.com/deepanshutr/lgtv-core).
It exposes two surfaces:
- **CLI** (`lgtv ...`) — built on `spf13/cobra`.
- **Telegram bot** (`lgtv tg-bot`) — built on `go-telegram/bot`.

Both call the same internal `core.Client` (which is just an HTTP client
against `lgtv-core`).

## Why combined in one repo?

The CLI and the bot share:
- the `core.Client` type and its methods
- config loading (`LGTV_CORE_URL`, `LGTV_AUTO_WAKE`)
- key-name → endpoint mapping

Splitting them would mean either a shared lib repo (overkill) or
duplicated logic (bad). Single repo, single binary, two `cobra` subtrees.

## Layout

```
cmd/lgtv/main.go              entrypoint, wires cobra root
internal/core/client.go       HTTP client for lgtv-core
internal/cli/                 cobra commands for `lgtv ...`
internal/tg/bot.go            Telegram bot (`lgtv tg-bot`)
internal/config/config.go     env-driven config (envconfig-style)
```

## Design rules

1. **No TV identifiers in this repo.** All TV-specific config (IP, MAC,
   client-key) lives in `lgtv-core`. This binary only knows about the
   *core daemon's* URL.
2. **One HTTP client, shared.** Both the CLI and the bot use the exact
   same `core.Client.DoX()` methods. Surface code only formats input and
   output.
3. **Bot is allow-listed by default.** No one outside `LGTV_TG_ALLOWED_USER_IDS`
   can send commands. Failing closed is the design.
4. **Auto-wake is opt-in.** The naive thing is to wake before every
   command. The right thing is to make it explicit so users can
   distinguish "TV is asleep" from "TV is unreachable".

## Gotchas

- `go-telegram/bot` swallows errors silently without `bot.WithDebug()`.
  Always enable in dev. (Same lesson learned in `orchctl`.)
- `cobra` flag inheritance: persistent flags on the root carry to subcommands,
  but the `lgtv tg-bot` subcommand uses *its own* config struct — don't be
  tempted to share `Cmd.Flags()` between cli and bot subtrees.
- Telegram allow-list parsing: input is a comma-separated string of
  decimal user IDs. Tabs/spaces are tolerated. Empty list = bot refuses
  everyone (failing closed).

## Building locally

```bash
go build ./cmd/lgtv
./lgtv state
```

## Testing

```bash
go test ./...
```

The HTTP client is tested against `httptest.Server`. The Telegram bot has
unit tests for the command-dispatch table; the actual `bot.New(...)` is
not exercised in unit tests.

## CI/CD

- `.github/workflows/ci.yml` — `go vet`, `staticcheck`, `go test ./...` on every PR.
- `.github/workflows/release.yml` — on tag, `goreleaser` builds linux+darwin
  binaries for amd64+arm64 and attaches to release.
- `gitleaks` runs on every push.

## Related

- `lgtv-core` — the Python daemon this CLI/bot talks to.
- `lgtv-mcp` — the MCP server; also a `lgtv-core` consumer, also written in Go.
