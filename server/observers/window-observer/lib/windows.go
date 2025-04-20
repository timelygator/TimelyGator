//go:build windows

package lib

import (
    "fmt"
    "path/filepath"
    "syscall"
    "unsafe"

    "golang.org/x/sys/windows"
)

var (
    user32                     = windows.NewLazySystemDLL("user32.dll")
    kernel32                   = windows.NewLazySystemDLL("kernel32.dll")
    procGetForegroundWindow    = user32.NewProc("GetForegroundWindow")
    procGetWindowTextW         = user32.NewProc("GetWindowTextW")
    procGetWindowTextLengthW   = user32.NewProc("GetWindowTextLengthW")
    procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
    procOpenProcess            = kernel32.NewProc("OpenProcess")
    procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
)

// GetActiveWindowHandle returns the HWND of the currently active window.
func GetActiveWindowHandle() windows.HWND {
    hwnd, _, _ := procGetForegroundWindow.Call()
    return windows.HWND(hwnd)
}

// GetWindowTitle returns the text/title of the given window handle.
func GetWindowTitle(hwnd windows.HWND) (string, error) {
    length, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
    if length == 0 {
        return "", nil
    }
    buf := make([]uint16, length+1)
    procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), length+1)
    return syscall.UTF16ToString(buf), nil
}

// GetAppPath returns the full executable path for the process owning hwnd.
func GetAppPath(hwnd windows.HWND) (string, error) {
    // Get PID
    var pid uint32
    procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

    // Open process with QUERY_LIMITED_INFORMATION
    const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
    handle, _, err := procOpenProcess.Call(
        uintptr(PROCESS_QUERY_LIMITED_INFORMATION),
        0, // inherit handle
        uintptr(pid),
    )
    if handle == 0 {
        return "", fmt.Errorf("OpenProcess failed: %v", err)
    }
    defer windows.CloseHandle(windows.Handle(handle))

    // Query full image name
    buf := make([]uint16, windows.MAX_PATH)
    size := uint32(len(buf))
    ret, _, err := procQueryFullProcessImageNameW.Call(
        handle,
        0,
        uintptr(unsafe.Pointer(&buf[0])),
        uintptr(unsafe.Pointer(&size)),
    )
    if ret == 0 {
        return "", fmt.Errorf("QueryFullProcessImageNameW failed: %v", err)
    }
    return syscall.UTF16ToString(buf[:size]), nil
}

// GetAppName returns the basename of the executable for hwnd.
func GetAppName(hwnd windows.HWND) (string, error) {
    path, err := GetAppPath(hwnd)
    if err != nil {
        return "", err
    }
    if path == "" {
        return "", nil
    }
    return filepath.Base(path), nil
}
