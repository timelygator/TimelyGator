package core

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"
)

// Merge merges src into dst, recursively. Values in src take precedence.
func Merge(dst, src map[string]interface{}) map[string]interface{} {
    for k, srcVal := range src {
        if dstVal, exists := dst[k]; exists {
            dstMap, dstIsMap := dstVal.(map[string]interface{})
            srcMap, srcIsMap := srcVal.(map[string]interface{})
            // If both are maps, merge recursively.
            if dstIsMap && srcIsMap {
                dst[k] = Merge(dstMap, srcMap)
                continue
            }
        }
        // Otherwise, take srcVal
        dst[k] = srcVal
    }
    return dst
}

// commentOutYAML prepends '#' to each non-empty line to "comment out" the YAML content.
func commentOutYAML(s string) string {
    lines := strings.Split(s, "\n")
    for i, line := range lines {
        if strings.TrimSpace(line) != "" {
            lines[i] = "#" + line
        }
    }
    return strings.Join(lines, "\n")
}

// LoadConfigYAML reads (and merges) YAML configuration from the defaultConfig
// string plus an existing config file (if present). If the config file doesn't
// exist yet, it writes out a commented version of defaultConfig.
func LoadConfigYAML(appName string, defaultConfig string) (map[string]interface{}, error) {
    // Parse the defaultConfig into a map
    defaultMap := make(map[string]interface{})
    if err := yaml.Unmarshal([]byte(defaultConfig), &defaultMap); err != nil {
        return nil, fmt.Errorf("default config is invalid YAML: %w", err)
    }

    // Determine the directory and ensure it exists
    configDir := GetConfigDir(appName)
    EnsurePathExists(configDir)

    // Build the path to the actual config file, e.g. ~/.config/tg-server/config.yaml
    configFilePath := filepath.Join(configDir, "config.yaml")

    // If file doesn't exist, create it with commented-out defaults
    if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
        log.Printf("Config file %s not found; creating with commented-out defaults.", configFilePath)
        if writeErr := os.WriteFile(
            configFilePath,
            []byte(commentOutYAML(defaultConfig)),
            0o644,
        ); writeErr != nil {
            return nil, fmt.Errorf("could not create config file: %w", writeErr)
        }
        return defaultMap, nil
    }

    // If it exists, read and parse
    data, err := os.ReadFile(configFilePath)
    if err != nil {
        return nil, fmt.Errorf("could not read config file %s: %w", configFilePath, err)
    }

    userMap := make(map[string]interface{})
    if err := yaml.Unmarshal(data, &userMap); err != nil {
        return nil, fmt.Errorf("config file is invalid YAML: %w", err)
    }

    // Merge user config into defaults
    finalMap := Merge(defaultMap, userMap)
    return finalMap, nil
}

// SaveConfigYAML writes the given config (as a YAML string) to
// "~/.config/<appName>/config.yaml". It validates the YAML before writing.
func SaveConfigYAML(appName, configString string) error {
    // Ensure it's valid YAML
    testMap := make(map[string]interface{})
    if err := yaml.Unmarshal([]byte(configString), &testMap); err != nil {
        return fmt.Errorf("config is invalid YAML: %w", err)
    }

    // Ensure the directory exists
    configDir := GetConfigDir(appName)
    EnsurePathExists(configDir)

    // Construct the file path
    configFilePath := filepath.Join(configDir, "config.yaml")

    // Write the YAML file
    if err := os.WriteFile(configFilePath, []byte(configString), 0o644); err != nil {
        return fmt.Errorf("could not write config file: %w", err)
    }
    return nil
}
