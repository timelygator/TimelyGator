//go:build !darwin

package lib

// GetCurrentWindowMacOS is a stub used when building on non‑macOS systems.
func GetCurrentWindowMacOS(strategy string) (*WindowInfo, error) {
    return nil, FatalError{"darwin build tag not enabled"}
}
