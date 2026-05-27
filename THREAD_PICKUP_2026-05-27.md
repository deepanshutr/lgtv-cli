# Thread pickup — 2026-05-27

## What shipped

One commit, pushed to `origin/main`:

- `e8edd03` — Add lgtv-display script for HDMI 4 input toggling

## Context

deePC is now physically wired into the LG OLED55CXPTA via HDMI 4 (see `~/.claude/projects/-home-deepanshutr/memory/project_depc_hdmi4_lgtv.md`). The TV doubles as the deePC's primary 4K display while remaining LAN-controllable via the lgtv stack.

The script gives a quick CLI toggle for what the TV is showing without touching xrandr — single HDMI output on this box, so mirror/extend modes are N/A.

## What's in the script

`scripts/lgtv-display`, also installed at `~/.local/bin/lgtv-display`:

| Subcommand | Action |
|---|---|
| `lgtv-display state` | pretty-prints `/state` as JSON |
| `lgtv-display tv-show` | wakes TV if off, then `POST /app/launch {"id": "com.webos.app.hdmi4"}` (deePC desktop on TV) |
| `lgtv-display tv-hide` | `POST /app/launch {"id": "youtube.leanback.v4"}` — TV switches to YouTube; deePC keeps rendering underneath |
| `lgtv-display tv-off` | `POST /power/off` — graceful standby; xrandr untouched |
| `--help` / `-h` | usage |

Smoke-tested live: tv-show / tv-hide round-trip confirmed against the running TV.

## Schema gotchas the script encodes (record these)

- Daemon's app-launch field is `id`, NOT `app_id`. (Both `/app/launch` and `/app/close` share the same `AppReq` schema.)
- Power-off endpoint is `POST /power/off`, NOT `/power_off`.
- `/input/switch` returns 502 from the TV when launching HDMI from a streaming-app context. Script uses `/app/launch` with the HDMI app id (`com.webos.app.hdmi4`) — works reliably across foreground states.

## Resume incantation

```bash
# the script is self-contained, no build:
which lgtv-display          # → ~/.local/bin/lgtv-display
lgtv-display state          # sanity
lgtv-display tv-hide        # TV → YouTube
lgtv-display tv-show        # TV → deePC desktop
```

## Open follow-ups

- Bind `Super+T` to toggle `tv-show` ↔ `tv-hide` via xbindkeys / GNOME settings.
- The Go `lgtv tg-bot` binary in this repo is NOT running anywhere (token isn't wired). orchctl-v2 owns the Telegram surface; this binary is reference code only as of 2026-05-13.
- HDMI-CEC is NOT viable on this box — NVIDIA GTX 1070 exposes no `/dev/cec*`. See `~/.claude/projects/-home-deepanshutr/memory/feedback_nvidia_no_cec.md`. LAN-only via `lgtv-core` covers everything CEC would have.
