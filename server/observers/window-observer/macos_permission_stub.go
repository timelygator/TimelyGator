//go:build !darwin
 
package windowobserver

// BackgroundEnsurePermissions is a no‑op on non‑macOS platforms.

func BackgroundEnsurePermissions() {

}