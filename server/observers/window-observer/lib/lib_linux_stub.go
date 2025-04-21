//go:build !linux

package lib

// GetCurrentWindowMacOS is a stub used when building on non‑macOS systems.
func GetCurrentWindowLinux() (*WindowInfo, error) {
	return nil, FatalError{"darwin build tag not enabled"}
}
