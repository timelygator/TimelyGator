package afkobserver

import (
    "fmt"
    "log"
    "time"

    "github.com/caarlos0/env"
    "github.com/joho/godotenv"
)

// AFKConfig holds only the relevant fields for AFK detection
type AFKConfig struct {
    Timeout  int `env:"TIMEOUT"   envDefault:"6"` // seconds
    PollTime int `env:"POLL_TIME" envDefault:"5"`   // seconds
}

// LoadAFKConfig loads .env (if present) and parses environment variables
// into an AFKConfig struct.
func LoadAFKConfig() (AFKConfig, error) {
    // 1) Attempt to load the .env file
    err := godotenv.Load()
    if err != nil {
        log.Printf("[LoadAFKConfig] Could not load .env file: %v (continuing)", err)
    }

    // 2) Parse environment variables into AFKConfig
    cfg := AFKConfig{}
    if err := env.Parse(&cfg); err != nil {
        return AFKConfig{}, fmt.Errorf("failed to parse environment: %w", err)
    }

    // Optional sanity check: ensure Timeout >= PollTime
    if cfg.Timeout < cfg.PollTime {
        log.Printf("[LoadAFKConfig] Warning: TIMEOUT (%d) < POLL_TIME (%d), which may cause issues.", cfg.Timeout, cfg.PollTime)
    }

    return cfg, nil
}

// Convert these int (seconds) to time.Duration for convenience.
func (a AFKConfig) TimeoutDuration() time.Duration {
    return time.Duration(a.Timeout) * time.Second
}

func (a AFKConfig) PollTimeDuration() time.Duration {
    return time.Duration(a.PollTime) * time.Second
}
