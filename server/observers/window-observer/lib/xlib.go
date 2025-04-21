//go:build linux

package lib

import (
	"errors"
	"log"
	"strings"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

// X is the shared X11 connection handle
var X *xgbutil.XUtil

// RootWin is the root window of the default screen
var RootWin xproto.Window

// Atoms used for window properties
var (
	AtomActiveWindow xproto.Atom
	AtomNETWMName    xproto.Atom
	AtomUTF8String   xproto.Atom
	AtomWMName       xproto.Atom
	AtomWMClass      xproto.Atom
	AtomNETWMPID     xproto.Atom
)

// FatalError is returned when the X connection is closed
var ErrFatal = errors.New("x server connection closed")

func init() {
	var err error
	X, err = xgbutil.NewConn()
	if err != nil {
		log.Fatalf("unable to connect to X server: %v", err)
	}
	RootWin = X.RootWin()

	AtomActiveWindow, _ = xprop.Atom(X, "_NET_ACTIVE_WINDOW", false)
	AtomNETWMName, _ = xprop.Atom(X, "_NET_WM_NAME", false)
	AtomUTF8String, _ = xprop.Atom(X, "UTF8_STRING", false)
	AtomWMName, _ = xprop.Atom(X, "WM_NAME", false)
	AtomWMClass, _ = xprop.Atom(X, "WM_CLASS", false)
	AtomNETWMPID, _ = xprop.Atom(X, "_NET_WM_PID", false)
}

// GetCurrentWindowID returns the XID of the currently active window, or 0 if none.
func GetCurrentWindowID() (xproto.Window, error) {
	reply, err := xproto.GetProperty(
		X.Conn(), false, RootWin,
		AtomActiveWindow, xproto.AtomAny,
		0, 1,
	).Reply()
	if err != nil {
		return 0, ErrFatal
	}
	if len(reply.Value) < 4 {
		return 0, nil
	}
	winID := xgb.Get32(reply.Value)
	if winID == 0 {
		return 0, nil
	}
	return xproto.Window(winID), nil
}

// GetWindowName returns the UTF-8 name of the given window, or "unknown" if unavailable.
func GetWindowName(win xproto.Window) string {
	name, err := xprop.PropValStr(xprop.GetProperty(X, win, "_NET_WM_NAME"))
	if err == nil && name != "" {
		return name
	}
	// Fallback to WM_NAME
	name, err = xprop.PropValStr(xprop.GetProperty(X, win, "WM_NAME"))
	if err == nil && name != "" {
		return name
	}
	return "unknown"
}

// GetWindowClass returns the class (window instance) of the given window, or "unknown".
func GetWindowClass(win xproto.Window) string {
	raw, err := xprop.PropValStr(xprop.GetProperty(X, win, "WM_CLASS"))
	if err != nil || raw == "" {
		return "unknown"
	}
	parts := strings.Split(raw, "\x00")
	if len(parts) >= 2 {
		return parts[1]
	}
	return parts[0]
}

// GetWindowPID returns the process ID owning the window, or an error.
func GetWindowPID(win xproto.Window) (uint32, error) {
	reply, err := xproto.GetProperty(
		X.Conn(), false, win,
		AtomNETWMPID, xproto.AtomAny,
		0, 1,
	).Reply()
	if err != nil {
		return 0, ErrFatal
	}
	if len(reply.Value) < 4 {
		return 0, errors.New("pid property not found")
	}
	return xgb.Get32(reply.Value), nil
}
