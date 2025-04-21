//go:build windows

package lib

func GetCurrentWindowWindows() (*WindowInfo, error) {
	hwnd := GetActiveWindowHandle()
	app, _ := GetAppName(hwnd)
	title, _ := GetWindowTitle(hwnd)

	if app == "" {
		app = "unknown"
	}
	if title == "" {
		title = "unknown"
	}
	return &WindowInfo{App: app, Title: title}, nil
}
