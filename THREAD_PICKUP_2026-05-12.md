# Thread pickup — 2026-05-12

## What shipped this thread

Initial scaffold of `lgtv-cli`: a Go Cobra CLI + Telegram bot subcommand
that talks HTTP to the local `lgtv-core` daemon (see sibling repo).

- `cmd/lgtv` — entrypoint
- `internal/cli` — Cobra command tree (`lgtv state|wake|power|vol|mute|app|key`)
- `internal/tg` — Telegram bot dispatcher (`lgtv tg-bot`) using `go-telegram/bot`,
  allow-listed by Telegram user ID
- `internal/core` — HTTP client for `lgtv-core` at `127.0.0.1:8765`
- 4 unit tests pass; `go vet` clean
- GitHub Actions CI (go vet, staticcheck, tests) + goreleaser release
  pipeline + gitleaks + dependabot
- MIT licensed, public

## Live state

- Binary built at `~/.local/bin/lgtv` — exercises end-to-end against the
  live `lgtv-core` daemon successfully (`lgtv state` returns JSON snapshot)
- **Telegram bot subcommand is NOT running** — the Telegram surface lives
  in `orchctl-v2` instead (operator decision: avoid two bot processes
  competing on the same token). This binary's `tg-bot` code is functional
  but token-less by design.

## Resume incantation

```bash
cd ~/github.com/deepanshutr/lgtv-cli
go vet ./... && go test ./... -count=1
go build -o ~/.local/bin/lgtv ./cmd/lgtv
~/.local/bin/lgtv state    # hits lgtv-core, returns JSON
```

To wire the standalone tg-bot (if ever needed):
```bash
export LGTV_TG_BOT_TOKEN=<from @BotFather>
export LGTV_TG_ALLOWED_USER_IDS=<your_tg_uid>
lgtv tg-bot
```

## Not done

- No standalone systemd unit (intentional — operator owns the bot in orchctl-v2)
- README documents the bot mode but it's effectively unwired in production

## Related repos

- `lgtv-core` — Python daemon this CLI calls
- `lgtv-mcp` — sibling MCP server, independent of this repo
- `orchctl-v2` — where the production `/tv` Telegram surface actually lives
