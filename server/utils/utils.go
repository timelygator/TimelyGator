package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"timelygator/server/database/models"
	"timelygator/server/utils/types"

	"github.com/adrg/xdg"
	"gorm.io/datatypes"
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

// mapToJSON converts a map into datatypes.JSON for storing in GORM.
func MapToJSON(data map[string]interface{}) (datatypes.JSON, error) {
	if data == nil {
		return datatypes.JSON([]byte("{}")), nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IsNotFound(err error) bool {
	// Return true if it matches
	return false
}

// parseIso8601 tries time.Parse with RFC3339 or similar
func ParseIso8601(val string) (time.Time, error) {
	return time.Parse(time.RFC3339, val)
}

// writeAttachmentJSON sets the Content-Disposition for file download
func WriteAttachmentJSON(w http.ResponseWriter, data interface{}, filename string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	json.NewEncoder(w).Encode(data)
}

// convertToEvent replicates the Python "Event(**data)" constructor
func ConvertToEvent(m map[string]interface{}) (*models.Event, error) {
	evt := &models.Event{}

	// 1) Parse ID (if any)
	if idVal, ok := m["id"].(float64); ok {
		evt.ID = uint(idVal) // numeric JSON => float64 => cast to uint
	}

	// 2) Parse timestamp (RFC3339 string)
	if tsVal, ok := m["timestamp"].(string); ok {
		t, err := ParseIso8601(tsVal)
		if err != nil {
			return nil, err
		}
		evt.Timestamp = t
	}

	// 3) Parse duration (float64 => time.Duration in seconds)
	if durVal, ok := m["duration"].(float64); ok {
		evt.Duration = time.Duration(durVal) * time.Second
	}

	// 4) Parse data (map[string]interface{}) => datatypes.JSON
	if dataVal, ok := m["data"].(map[string]interface{}); ok {
		b, err := json.Marshal(dataVal)
		if err != nil {
			return nil, err
		}
		// Set evt.Data to the raw JSON bytes
		evt.Data = datatypes.JSON(b)
	}

	return evt, nil
}
