// Package cli wires up the cobra command tree for `lgtv`.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deepanshutr/lgtv-cli/internal/config"
	"github.com/deepanshutr/lgtv-cli/internal/core"
	"github.com/deepanshutr/lgtv-cli/internal/tg"
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	cfg := config.Load()
	client := core.New(cfg.CoreURL)
	ctx := context.Background()

	root := &cobra.Command{
		Use:   "lgtv",
		Short: "Remote-control your LG webOS TV (via local lgtv-core daemon)",
	}
	root.PersistentFlags().BoolVar(&cfg.AutoWake, "auto-wake", cfg.AutoWake,
		"Send WoL + wait for TV before each command")

	wake := func() error {
		if !cfg.AutoWake {
			return nil
		}
		return client.Wake(ctx)
	}

	root.AddCommand(
		&cobra.Command{
			Use:   "wake",
			Short: "Send a wake-on-LAN magic packet and wait for the TV",
			RunE: func(_ *cobra.Command, _ []string) error {
				return client.Wake(ctx)
			},
		},
		&cobra.Command{
			Use:   "power [off]",
			Short: "Power the TV on (via wake) or off",
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) > 0 && args[0] == "off" {
					return client.PowerOff(ctx)
				}
				return client.Wake(ctx)
			},
		},
		&cobra.Command{
			Use:   "state",
			Short: "Show current TV state (volume, app, etc.)",
			RunE: func(_ *cobra.Command, _ []string) error {
				if err := wake(); err != nil {
					return err
				}
				s, err := client.State(ctx)
				if err != nil {
					return err
				}
				b, _ := json.MarshalIndent(s, "", "  ")
				fmt.Println(string(b))
				return nil
			},
		},
		newVolumeCmd(client, wake),
		&cobra.Command{
			Use:   "mute [on|off|toggle]",
			Short: "Mute control",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := wake(); err != nil {
					return err
				}
				switch strings.ToLower(args[0]) {
				case "on", "true", "1":
					return client.Mute(ctx, true)
				case "off", "false", "0":
					return client.Mute(ctx, false)
				default:
					return fmt.Errorf("expected on|off, got %s", args[0])
				}
			},
		},
		newAppCmd(client, wake),
		&cobra.Command{
			Use:   "key <name>",
			Short: "Press a remote key (HOME, BACK, UP, DOWN, ...)",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := wake(); err != nil {
					return err
				}
				return client.PressKey(ctx, args[0])
			},
		},
		newTGBotCmd(cfg, client),
	)
	return root
}

func newVolumeCmd(client *core.Client, wake func() error) *cobra.Command {
	return &cobra.Command{
		Use:     "vol <level | +N | -N>",
		Aliases: []string{"volume"},
		Short:   "Set absolute volume (0-100) or apply a relative delta (+N / -N)",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := wake(); err != nil {
				return err
			}
			ctx := context.Background()
			v := args[0]
			if strings.HasPrefix(v, "+") || strings.HasPrefix(v, "-") {
				delta, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				return client.VolumeDelta(ctx, delta)
			}
			level, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			return client.VolumeAbsolute(ctx, level)
		},
	}
}

func newAppCmd(client *core.Client, wake func() error) *cobra.Command {
	c := &cobra.Command{Use: "app", Short: "App control"}
	c.AddCommand(&cobra.Command{
		Use:   "launch <app-id>",
		Short: "Launch an app by ID (e.g. netflix, youtube.leanback.v4)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := wake(); err != nil {
				return err
			}
			return client.LaunchApp(context.Background(), args[0])
		},
	})
	return c
}

func newTGBotCmd(cfg config.Config, client *core.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "tg-bot",
		Short: "Run the Telegram bot (long-poll, allow-listed users only)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.TGToken == "" {
				return fmt.Errorf("set LGTV_TG_BOT_TOKEN to run the bot")
			}
			if len(cfg.AllowedIDs) == 0 {
				fmt.Fprintln(os.Stderr,
					"warning: LGTV_TG_ALLOWED_USER_IDS is empty; bot will reject all users")
			}
			return tg.Run(cmd.Context(), cfg.TGToken, cfg.AllowedIDs, client)
		},
	}
}
