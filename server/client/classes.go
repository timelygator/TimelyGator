package client

import (
	"encoding/json"
	"log"
)

// CategorySpec is "spec" map (e.g. {"type":"regex","regex":"...","ignore_case":true})
type CategorySpec map[string]interface{}

// ClassItem represents a classification entry:
//   - Name => slice of strings (like ["Work","Programming"])
//   - Rule => a CategorySpec containing "type", "regex", etc.
type ClassItem struct {
	Name []string     `json:"name"`
	Rule CategorySpec `json:"rule"`
}

// DefaultClasses is the fallback if the server does not provide "classes" or if parsing fails.
var DefaultClasses = []ClassItem{
	{
		Name: []string{"Work"},
		Rule: CategorySpec{
			"type":  "regex",
			"regex": "Google Docs|libreoffice|ReText",
		},
	},
	{
		Name: []string{"Work", "Programming"},
		Rule: CategorySpec{
			"type":  "regex",
			"regex": "|VS Code|GitHub|Stack Overflow|BitBucket|Gitlab|vim",
		},
	},
	{
		Name: []string{"Work", "Programming", "TimelyGator"},
		Rule: CategorySpec{
			"type":        "regex",
			"regex":       "TimelyGator|tg-",
			"ignore_case": true,
		},
	},
	{
		Name: []string{"Media", "Social Media"},
		Rule: CategorySpec{
			"type":        "regex",
			"regex":       "reddit|Facebook|Twitter|Instagram|devRant",
			"ignore_case": true,
		},
	},
	{
		Name: []string{"Media", "Music"},
		Rule: CategorySpec{
			"type":        "regex",
			"regex":       "Spotify|Deezer",
			"ignore_case": true,
		},
	},
	{
		Name: []string{"Comms", "IM"},
		Rule: CategorySpec{
			"type":  "regex",
			"regex": "Messenger|Telegram|Signal|WhatsApp|Slack|Discord",
		},
	},
	{
		Name: []string{"Comms", "Email"},
		Rule: CategorySpec{
			"type":  "regex",
			"regex": "Gmail|Thunderbird",
		},
	},
}

// GetClasses attempts to fetch "classes" from the server
// via c.GetSetting("classes"). If that fails (error or parse issue),
// it returns defaultClasses as a fallback.
func GetClasses(c *TimelyGatorClient) []ClassItem {
	// Attempt to fetch from the server
	settingKey := "classes"
	data, err := c.GetSetting(&settingKey)
	if err != nil {
		log.Printf("Failed to get classes from server; using default: %v", err)
		return DefaultClasses
	}
	if data == nil {
		log.Println("Received nil data for classes; using default classes.")
		return DefaultClasses
	}

	// The server is expected to return an array of objects shaped like:
	//  [
	//    { "name": ["Work"], "rule": { "type":"regex","regex":"somepattern" } },
	//    ...
	//  ]
	// We unmarshal that into []ClassItem
	rawJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal 'classes' setting to JSON; default used: %v", err)
		return DefaultClasses
	}

	var classesFromServer []ClassItem
	if err := json.Unmarshal(rawJSON, &classesFromServer); err != nil {
		log.Printf("Failed to unmarshal 'classes' data; default used: %v", err)
		return DefaultClasses
	}

	return classesFromServer
}
