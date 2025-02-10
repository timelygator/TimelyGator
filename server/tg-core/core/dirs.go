package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

func EnsurePathExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

func EnsureReturnedPathExists(f func(string) string) func(string) string {
	return func(subpath string) string {
		path := f(subpath)
		EnsurePathExists(path)
		return path
	}
}

func GetDataDir(moduleName string) string {
	dataDir := xdg.DataHome
	if moduleName != "" {
		return filepath.Join(dataDir, moduleName)
	}
	return dataDir
}

func GetCacheDir(moduleName string) string {
	cacheDir := xdg.CacheHome
	if moduleName != "" {
		return filepath.Join(cacheDir, moduleName)
	}
	return cacheDir
}

func GetConfigDir(moduleName string) string {
	configDir := xdg.ConfigHome
	if moduleName != "" {
		return filepath.Join(configDir, moduleName)
	}
	return configDir
}

func GetLogDir(moduleName string) string {
	var logDir string
	if runtime.GOOS == "linux" {
		logDir = filepath.Join(xdg.CacheHome, "log")
	} else {
		logDir = xdg.StateHome
	}
	if moduleName != "" {
		return filepath.Join(logDir, moduleName)
	}
	return logDir
}

func main() {
	fmt.Println("Data Dir:", GetDataDir("timelygator"))
	fmt.Println("Cache Dir:", GetCacheDir("timelygator"))
	fmt.Println("Config Dir:", GetConfigDir("timelygator"))
	fmt.Println("Log Dir:", GetLogDir("timelygator"))
}
