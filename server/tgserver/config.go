package tgserver

import (
	"fmt"

	"timelygator/server/tg-core/core"
)

var defaultConfigYAML = `
server:
  host: "localhost"
  port: "5600"
  storage: "sqlite"
  cors_origins: ""
  custom_static: {}

server-testing:
  host: "localhost"
  port: "5666"
  storage: "sqlite"
  cors_origins: ""
  custom_static: {}
`

// Config is a map that holds the entire parsed YAML.
// Typically: map[string]interface{} with keys like "server", "server-testing".
var Config map[string]interface{}

func init() {
	result, err := core.LoadConfigYAML("tg-server", defaultConfigYAML)
	if err != nil {
		panic(fmt.Errorf("failed to load config YAML: %w", err))
	}
	Config = result
}

func FromSection(name string) (map[string]interface{}, error) {
	sectionVal, ok := Config[name]
	if !ok {
		return nil, fmt.Errorf("no section %q found in config", name)
	}
	sectionMap, ok := sectionVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("section %q is not a map in config", name)
	}
	return sectionMap, nil
}
