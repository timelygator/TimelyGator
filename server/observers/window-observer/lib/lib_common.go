//go:build linux || darwin || windows

package lib

import (
    "fmt"
    "runtime"
)

// ---------- shared types ----------------------------------------------------

type WindowInfo struct {
    App       string  `json:"app"`
    Title     string  `json:"title"`
    URL       *string `json:"url,omitempty"`
    Incognito *bool   `json:"incognito,omitempty"`
}

type FatalError struct{ msg string }

func (f FatalError) Error() string { return f.msg }

// ---------- public dispatcher ----------------------------------------------

func GetCurrentWindow(strategy string) (*WindowInfo, error) {
    switch runtime.GOOS {
    case "linux":
        return GetCurrentWindowLinux()
    case "windows":
        return GetCurrentWindowWindows()
    case "darwin":
        if strategy == "" {
            return nil, FatalError{"macOS strategy not specified"}
        }
        return GetCurrentWindowMacOS(strategy)
    default:
        return nil, FatalError{fmt.Sprintf("unknown platform: %s", runtime.GOOS)}
    }
}

// ---------- helpers --------------------------------------------------------

func StringPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

func boolPtr(b bool) *bool {
    return &b
}
