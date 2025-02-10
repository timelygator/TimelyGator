package types

const ModuleName = "timelygator"

type Config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"development"`
	Domain             string `env:"DOMAIN" envDefault:"localhost"`
	Port               string `env:"PORT" envDefault:"8080"`
	DataSourceName     string `env:"DSN" envDefault:"timelygator.db"` // SQLite - file.db, MySQL - user:password@tcp(localhost:3306)/dbname
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}
