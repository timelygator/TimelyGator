package windowobserver

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

//   HOST             → server host (string)
//   PORT             → server port (int, default 8080)
//   TESTING          → enable testing mode (bool)
//   VERBOSE          → verbose logs (bool)
//   EXCLUDE_TITLE    → exclude empty titles (bool)
//   EXCLUDE_TITLES   → comma‑separated list of regexp strings to ignore
//   POLL_TIME        → sampling interval in seconds (float, default 1.0)
//   STRATEGY         → macOS only: jxa | applescript | swift  (default swift)
//
// You can override any of these at runtime with CLI flags if desired; the
// observer’s flag parser should fall back to the values supplied here.
//
// Note: `envSeparator` parses `EXCLUDE_TITLES="Teams,Zoom"` into a slice.
// -----------------------------------------------------------------------------

type WindowObserverConfig struct {
	Host          string   `env:"HOST"`
	Port          string   `env:"PORT"         envDefault:"8080"`
	Testing       bool     `env:"TESTING"      envDefault:"false"`
	Verbose       bool     `env:"VERBOSE"      envDefault:"false"`
	ExcludeTitle  bool     `env:"EXCLUDE_TITLE" envDefault:"false"`
	ExcludeTitles []string `env:"EXCLUDE_TITLES" envSeparator:","`
	PollTime      float64  `env:"POLL_TIME"    envDefault:"1.0"`
	Strategy      string   `env:"STRATEGY"     envDefault:"swift"`
}

// LoadConfig reads .env (if present) and environment variables into the struct.
func LoadConfig() (WindowObserverConfig, error) {
	_ = godotenv.Load() // optional; ignore error if .env absent

	cfg := WindowObserverConfig{}
	if err := env.Parse(&cfg); err != nil {
		return WindowObserverConfig{}, fmt.Errorf("window‑observer config: %w", err)
	}

	// Normalise strategy to lowercase for easy comparison.
	cfg.Strategy = strings.ToLower(cfg.Strategy)
	switch cfg.Strategy {
	case "jxa", "applescript", "swift":
	default:
		log.Printf("[LoadConfig] unknown STRATEGY %q – falling back to 'swift'", cfg.Strategy)
		cfg.Strategy = "swift"
	}

	return cfg, nil
}

// PollDuration converts the float seconds into a time.Duration.
func (c WindowObserverConfig) PollDuration() time.Duration {
	return time.Duration(c.PollTime * float64(time.Second))
}
