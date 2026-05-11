# lgtv-cli

A small Go CLI and Telegram bot for controlling your LG webOS TV. Both are
thin shells over [`lgtv-core`](https://github.com/deepanshutr/lgtv-core),
which owns the actual WebSocket session with the TV.

```
┌─ lgtv-core ─────────────────────────┐
│  Python daemon                      │
│  127.0.0.1:8765                     │
└──────┬───────────────────────┬──────┘
       │                       │
       ▼                       ▼
 ┌──────────┐           ┌────────────────┐
 │ lgtv     │           │ lgtv tg-bot    │
 │ (CLI)    │           │ (Telegram)     │
 └──────────┘           └────────────────┘
```

## Install

```bash
go install github.com/deepanshutr/lgtv-cli/cmd/lgtv@latest
```

Or grab a release binary from the [releases page](https://github.com/deepanshutr/lgtv-cli/releases).

## Quick start

Assumes `lgtv-core` is already running at `127.0.0.1:8765`.

```bash
lgtv power off
lgtv vol 20
lgtv vol +3
lgtv app launch netflix
lgtv app list
lgtv key home
lgtv state
```

All commands accept `--auto-wake` (or `LGTV_AUTO_WAKE=1`), which fires a
`/wake` first and waits for the TV to come up before sending the real
command.

## Telegram bot mode

```bash
export LGTV_TG_BOT_TOKEN=...           # from @BotFather
export LGTV_TG_ALLOWED_USER_IDS=12345,67890
lgtv tg-bot
```

The bot only responds to allow-listed Telegram user IDs. Commands:

```
/power_off       — standby
/wake            — wake from standby
/vol 30          — absolute volume
/vol +3          — relative volume
/mute            — toggle mute
/app netflix     — launch app
/state           — fetch current state
/key home        — press a remote key
```

## Configuration

| Env | Default | Meaning |
|-----|---------|---------|
| `LGTV_CORE_URL` | `http://127.0.0.1:8765` | Where `lgtv-core` is listening |
| `LGTV_AUTO_WAKE` | `false` | Wake the TV before each command |
| `LGTV_TG_BOT_TOKEN` | — | Telegram bot token (bot mode only) |
| `LGTV_TG_ALLOWED_USER_IDS` | — | Comma-separated allow-list of Telegram user IDs |

## License

MIT — see `LICENSE`.
