//go:build darwin

package lib

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
)

const appleScriptSource = `
global frontApp, frontAppName, windowTitle

set windowTitle to ""
tell application "System Events"
    set frontApp to first application process whose frontmost is true
    set frontAppName to name of frontApp
    tell process frontAppName
        try
            tell (1st window whose value of attribute "AXMain" is true)
                set windowTitle to value of attribute "AXTitle"
            end tell
        end try
    end tell
end tell

return frontAppName & "\n" & windowTitle
`

// GetInfo executes the embedded AppleScript and returns {"app": appName, "title": windowTitle}.
func GetInfoAppleScript() (map[string]string, error) {
    cmd := exec.Command("osascript", "-e", appleScriptSource)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("osascript error: %v | stderr: %s", err, stderr.String())
    }

    // osascript prints with a trailing newline; trim and split.
    output := strings.TrimSpace(stdout.String())
    parts := strings.SplitN(output, "\n", 2)
    if len(parts) != 2 {
        return nil, fmt.Errorf("unexpected output format: %q", output)
    }

    return map[string]string{
        "app":   parts[0],
        "title": parts[1],
    }, nil
}
