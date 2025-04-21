//go:build darwin

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

// getInfo calls the JXA script (printAppStatus.jxa) using osascript and parses its JSON output.
func GetInfo() (map[string]string, error) {
	scriptPath, err := filepath.Abs("lib/printAppStatus.jxa")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve script path: %w", err)
	}

	cmd := exec.Command("osascript", "-l", "JavaScript", scriptPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("osascript error: %v\nstderr: %s", err, stderr.String())
	}

	var result map[string]string
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JXA JSON: %w", err)
	}
	return result, nil
}
