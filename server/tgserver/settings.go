package tgserver

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"

	"timelygator/server/tg-core/core"
)

type Settings struct {
    configFile string
    data       map[string]interface{}
}

// NewSettings picks which file to use and then loads the contents.
func NewSettings(testing bool) (*Settings, error) {
    filename := "settings.json"
    if testing {
        filename = "settings-testing.json"
    }

    // Build path:
    configDir := core.GetConfigDir("tg-server")
    configFile := filepath.Join(configDir, filename)

    s := &Settings{
        configFile: configFile,
        data:       make(map[string]interface{}),
    }
    // Load existing data if the file exists
    if err := s.Load(); err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            return nil, err
        }
    }
    return s, nil
}

// Load reads the JSON file on disk into s.data.
func (s *Settings) Load() error {
    // Check if the file exists
    fi, err := os.Stat(s.configFile)
    if err != nil {
        return err
    }
    if fi.IsDir() {
        return fmt.Errorf("config file path is a directory: %s", s.configFile)
    }
    // Open and decode
    f, err := os.Open(s.configFile)
    if err != nil {
        return err
    }
    defer f.Close()

    decoder := json.NewDecoder(f)
    if err := decoder.Decode(&s.data); err != nil {
        return fmt.Errorf("failed to decode JSON settings: %w", err)
    }
    return nil
}

// Save writes s.data to disk as pretty-printed JSON.
func (s *Settings) Save() error {
    // Ensure parent dir exists
    if err := os.MkdirAll(filepath.Dir(s.configFile), 0755); err != nil {
        return fmt.Errorf("failed to create config dir: %w", err)
    }

    f, err := os.Create(s.configFile)
    if err != nil {
        return fmt.Errorf("failed to create config file: %w", err)
    }
    defer f.Close()

    encoder := json.NewEncoder(f)
    encoder.SetIndent("", "    ")
    if err := encoder.Encode(s.data); err != nil {
        return fmt.Errorf("failed to encode JSON settings: %w", err)
    }
    return nil
}

// Get retrieves a value by key. If key is empty, returns the entire map.
// If key is not found, returns defaultVal.
func (s *Settings) Get(key string, defaultVal interface{}) interface{} {
    if key == "" {
        return s.data
    }
    val, ok := s.data[key]
    if !ok {
        return defaultVal
    }
    return val
}

// Set sets (or deletes) a key/value in s.data and saves to disk.
// If value is nil, the key is removed.
func (s *Settings) Set(key string, value interface{}) error {
    if value != nil {
        s.data[key] = value
    } else {
        delete(s.data, key)
    }
    return s.Save()
}
