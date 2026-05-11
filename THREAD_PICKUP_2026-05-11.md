# Thread pickup — 2026-05-11 (lgtv-cli)

## What shipped this session

Initial scaffolding of the Go Cobra CLI + Telegram bot (`lgtv tg-bot`).
Single commit: `4d69742` Initial scaffolding.

- Cobra subcommands hit `lgtv-core`'s HTTP API at `127.0.0.1:8765`.
- `lgtv tg-bot` long-polls Telegram and forwards commands to the daemon
  (allow-listed by Telegram user ID).
- CI: gitleaks + `go test ./...`. Goreleaser config in `.goreleaser.yml`
  (multi-arch builds, dummy until first tag).
- Notify-on-release in `release.yml` gated on `TG_BOT_TOKEN` + `TG_CHAT_ID`
  GitHub Secrets — both unset, so notify is a no-op until set.

## Live state

- Binary symlinked at `~/.local/bin/lgtv`.
- CLI works end-to-end against the running `lgtv-core` daemon.

## Not yet wired

- **`tg-bot` has no token configured.** To enable:
  1. Create a bot via @BotFather.
  2. `LGTV_TG_BOT_TOKEN=… LGTV_TG_ALLOWED_USER_IDS=<your_tg_uid> lgtv tg-bot`.
  3. Write a systemd user unit (model on `lgtv-core.service`) so it persists.

## Gotchas to NOT re-discover

1. Local Go env: `~/.profile` exports `GOROOT=/home/deepanshutr/go/go1.18`
   which mismatches the brew go at `/home/linuxbrew/.linuxbrew/bin/go`
   (1.26.2). Always `unset GOROOT; export
   GOPROXY=https://proxy.golang.org,direct` before any `go` command.
2. Push email must be the noreply:
   `52166434+deepanshutr@users.noreply.github.com`. Per-repo only.
3. Telegram bot — only ONE process can poll per token. Before starting a
   new instance, `pgrep -fa lgtv | grep tg-bot` and kill orphans
   (`go-telegram/bot` swallows poll-conflict errors without `WithDebug`).

## Exact resume incantation

```bash
cd ~/github.com/deepanshutr/lgtv-cli
unset GOROOT; export GOPROXY=https://proxy.golang.org,direct
go build -o ~/.local/bin/lgtv .
lgtv state
lgtv volume -1
lgtv launch youtube.leanback.v4
# tg-bot (after setting env):
LGTV_TG_BOT_TOKEN=… LGTV_TG_ALLOWED_USER_IDS=… lgtv tg-bot
```

## Repo state at thread-close

- Branch: `main`, up to date with `origin/main`, clean tree.
- Single commit since init.
- CI green; goreleaser config in place but no tag yet.

## Memory references

- `project_lgtv_stack.md` — full project / sibling-repo overview
- `feedback_orchctl_debugging.md` — Telegram poll-conflict / go-telegram/bot
  silent-error patterns (apply identically here)
