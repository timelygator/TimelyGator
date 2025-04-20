//go:build linux

package lib

func GetCurrentWindowLinux() (*WindowInfo, error) {
	win, err := GetCurrentWindowID()
	if err != nil || win == 0 {
		return &WindowInfo{App: "unknown", Title: "unknown"}, err
	}
	app := GetWindowClass(win)
	title := GetWindowName(win)
	if app == "" {
		app = "unknown"
	}
	if title == "" {
		title = "unknown"
	}
	return &WindowInfo{App: app, Title: title}, nil
}
