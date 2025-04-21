//go:build windows
// +build windows

package lib

import "testing"

func TestGetCurrentWindowWindows(t *testing.T) {
	t.Skip("requires Windows environment or stubs for GetActiveWindowHandle, GetAppName, GetWindowTitle")
}
