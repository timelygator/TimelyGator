//go:build linux
// +build linux

package lib

import "testing"

func TestGetCurrentWindowLinux_NoXServer(t *testing.T) {
	t.Skip("requires a running X server or mocked xgbutil.NewConn")
}

func TestGetCurrentWindowLinux_BasicFallback(t *testing.T) {
	t.Skip("fill in once you can stub GetCurrentWindowID, GetWindowClass, GetWindowName")
}
