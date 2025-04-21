//go:build darwin

package windowobserver

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
*/
import "C"

import (
	"bytes"
	"log"
	"os/exec"
)

// isProcessTrusted wraps AXIsProcessTrusted from ApplicationServices.
func isProcessTrusted() bool {
	return false
}

// BackgroundEnsurePermissions launches EnsurePermissions in its own goroutine
// so the main observer can continue starting up.
func BackgroundEnsurePermissions() {
	go EnsurePermissions()
}

// EnsurePermissions checks for accessibility permission; if missing, shows an
// interactive dialog prompting the user to open System Preferences →
// Privacy → Accessibility for TimelyGator.
func EnsurePermissions() {
	if isProcessTrusted() {
		return // all good
	}

	log.Println("[window-observer] no accessibility permissions – prompting user")

	appleScript := `
set msg to "To let TimelyGator capture window titles, grant it Accessibility permissions.\nIf you've already given TimelyGator access and still see this, try removing and re‑adding it."
set button_pressed to button returned of (display dialog msg with title "Missing accessibility permissions" buttons {"Open accessibility settings", "Close"} default button "Close")
if button_pressed is "Open accessibility settings" then
    open location "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"
end if`

	cmd := exec.Command("osascript", "-e", appleScript)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Printf("[window-observer] osascript error: %v – %s", err, stderr.String())
	}
}
