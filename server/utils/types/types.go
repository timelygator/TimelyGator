package types

import (
	"fmt"

	"gorm.io/datatypes"
)

const ModuleName = "timelygator"
const ModuleVersion = "0.1.0"

type Config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"development"`
	Interface          string `env:"INTERFACE" envDefault:"localhost"`
	Port               string `env:"PORT" envDefault:"8080"`
	DataSourceName     string `env:"DSN" envDefault:"timelygator.db"` // SQLite - file.db, MySQL - user:password@tcp(localhost:3306)/dbname
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

type InfoResponse datatypes.JSON
type HTTPError string

// BucketCreationPayload is the payload for creating a bucket.
type BucketCreationPayload struct {
	Client   string `json:"client"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
}

// BucketUpdatePayload is the payload for updating a bucket.
type BucketUpdatePayload struct {
	Client   *string                `json:"client"`
	Type     *string                `json:"type"`
	Hostname *string                `json:"hostname"`
	Data     map[string]interface{} `json:"data"`
}

// ImportPayload is the payload for importing buckets.
type ImportPayload struct {
	Buckets map[string]interface{} `json:"buckets"`
}

// QueryPayload is the payload for querying events.
type QueryPayload struct {
	Timeperiods []string `json:"timeperiods"`
	Query       []string `json:"query"`
}

type NotFound struct {
	Code    string
	Message string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
