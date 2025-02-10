package utils

import (
	"errors"
	"os"
	"path/filepath"
	"timelygator/server/utils/types"

	"github.com/adrg/xdg"
)

func EnsurePathExists(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func GetDir(base string) (string, error) {
	path := ""
	switch base {
	case "data":
		path = filepath.Join(xdg.DataHome, types.ModuleName)
	case "cache":
		path = filepath.Join(xdg.CacheHome, types.ModuleName)
	case "config":
		path = filepath.Join(xdg.ConfigHome, types.ModuleName)
	case "log":
		path = filepath.Join(xdg.StateHome, types.ModuleName)
	}
	err := EnsurePathExists(path)
	return path, err
}
