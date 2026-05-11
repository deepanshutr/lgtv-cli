// Package tg implements the Telegram bot mode of lgtv.
package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/deepanshutr/lgtv-cli/internal/core"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Run starts a long-poll bot. Blocks until ctx is cancelled.
func Run(ctx context.Context, token string, allowed []int64, client *core.Client) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler(client, allowed)),
		bot.WithAllowedUpdates(bot.AllowedUpdates{"message"}),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return err
	}
	b.Start(ctx)
	return nil
}

func handler(client *core.Client, allowed []int64) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, u *models.Update) {
		if u.Message == nil {
			return
		}
		uid := u.Message.From.ID
		if !slices.Contains(allowed, uid) {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: u.Message.Chat.ID,
				Text:   "Not authorized.",
			})
			return
		}
		reply := dispatch(ctx, client, u.Message.Text)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: u.Message.Chat.ID,
			Text:   reply,
		})
	}
}

func dispatch(ctx context.Context, client *core.Client, text string) string {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "Empty command. Try /state, /vol 20, /app netflix, /key home."
	}
	cmd := strings.TrimPrefix(strings.ToLower(parts[0]), "/")
	args := parts[1:]

	switch cmd {
	case "wake":
		if err := client.Wake(ctx); err != nil {
			return errStr(err)
		}
		return "Awake."
	case "power_off", "off":
		if err := client.PowerOff(ctx); err != nil {
			return errStr(err)
		}
		return "TV off."
	case "state":
		s, err := client.State(ctx)
		if err != nil {
			return errStr(err)
		}
		b, _ := json.MarshalIndent(s, "", "  ")
		return "```\n" + string(b) + "\n```"
	case "vol", "volume":
		if len(args) != 1 {
			return "Usage: /vol <level | +N | -N>"
		}
		v := args[0]
		if strings.HasPrefix(v, "+") || strings.HasPrefix(v, "-") {
			d, err := strconv.Atoi(v)
			if err != nil {
				return errStr(err)
			}
			if err := client.VolumeDelta(ctx, d); err != nil {
				return errStr(err)
			}
			return fmt.Sprintf("Volume %+d", d)
		}
		l, err := strconv.Atoi(v)
		if err != nil {
			return errStr(err)
		}
		if err := client.VolumeAbsolute(ctx, l); err != nil {
			return errStr(err)
		}
		return fmt.Sprintf("Volume = %d", l)
	case "mute":
		on := len(args) == 0 || args[0] == "on"
		if err := client.Mute(ctx, on); err != nil {
			return errStr(err)
		}
		return fmt.Sprintf("Mute = %v", on)
	case "app":
		if len(args) != 1 {
			return "Usage: /app <app-id>"
		}
		if err := client.LaunchApp(ctx, args[0]); err != nil {
			return errStr(err)
		}
		return "Launched " + args[0]
	case "key":
		if len(args) != 1 {
			return "Usage: /key <name>"
		}
		if err := client.PressKey(ctx, args[0]); err != nil {
			return errStr(err)
		}
		return "Pressed " + args[0]
	default:
		return "Unknown command. Try: /state /wake /off /vol /mute /app /key"
	}
}

func errStr(err error) string {
	return "Error: " + err.Error()
}
