//go:build !windows

package lib

// GetCurrentWindowWindows is a stub used when building on non‑Windows systems.
func GetCurrentWindowWindows() (*WindowInfo, error) {
	return nil, FatalError{"windows build tag not enabled"}
}
