package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	CoreURL     string
	AutoWake    bool
	TGToken     string
	AllowedIDs  []int64
}

func Load() Config {
	c := Config{
		CoreURL:  envOr("LGTV_CORE_URL", "http://127.0.0.1:8765"),
		AutoWake: envBool("LGTV_AUTO_WAKE"),
		TGToken:  os.Getenv("LGTV_TG_BOT_TOKEN"),
	}
	for _, p := range strings.Split(os.Getenv("LGTV_TG_ALLOWED_USER_IDS"), ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if id, err := strconv.ParseInt(p, 10, 64); err == nil {
			c.AllowedIDs = append(c.AllowedIDs, id)
		}
	}
	return c
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func envBool(key string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}
