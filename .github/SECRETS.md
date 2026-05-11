# GitHub Secrets used by this repo

| Secret | Used in | Purpose |
|--------|---------|---------|
| `TG_BOT_TOKEN` | `release.yml` notify job | Post release notifications to operator's Telegram. Optional — notify is no-op if unset. |
| `TG_CHAT_ID`   | `release.yml` notify job | Target chat. |

The bot needs `LGTV_TG_BOT_TOKEN` and `LGTV_TG_ALLOWED_USER_IDS` *at runtime*
on the host that runs `lgtv tg-bot`. Those are **not** GitHub Actions secrets —
they live in your local systemd unit's `EnvironmentFile`.
