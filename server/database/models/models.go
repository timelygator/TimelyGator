package models

import (
	"encoding/json"
	"log"
	"time"

	"gorm.io/datatypes"
)

// For convenience, define aliases or helper types as needed.
// But we remove the original `type Data map[string]interface{}`
// in favor of GORM's datatypes.JSON.

type ConvertibleTimestamp interface{}
type Duration interface{}

// Event is stored in the DB with a JSON blob for Data.
type Event struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	BucketID  string         `gorm:"index" json:"bucket_id"`
	Timestamp time.Time      `gorm:"not null;type:timestamp" json:"timestamp"`
	Duration  time.Duration  `gorm:"not null" json:"duration"`
	Data      datatypes.JSON `gorm:"type:json" json:"data"`
}

// Bucket is also stored in the DB with a JSON blob for Data.
type Bucket struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Name     *string
	Type     string
	Client   string
	Hostname string
	Created  time.Time
	Data     datatypes.JSON `gorm:"type:json" json:"data"`
}

// NewEvent creates an Event with typed timestamp/duration
// and converts a map[string]interface{} (if any) into JSON.
func NewEvent(
	id uint,
	timestamp ConvertibleTimestamp,
	duration Duration,
	data map[string]interface{},
) *Event {

	// If no timestamp is provided, use "now".
	if timestamp == nil {
		log.Printf("Event initializer did not receive a timestamp argument, using now.")
		timestamp = time.Now().UTC()
	} else {
		timestamp = timestampParse(timestamp)
	}

	// Convert the map data into JSON (if any).
	var jsonData datatypes.JSON
	if data == nil {
		jsonData = datatypes.JSON([]byte("{}"))
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Failed to marshal event data: %v", err)
		}
		jsonData = b
	}

	return &Event{
		ID:        id,
		Timestamp: timestamp.(time.Time),
		Duration:  parseDuration(duration),
		Data:      jsonData,
	}
}

// timestampParse attempts to parse an RFC3339 or accept time.Time.
func timestampParse(tsIn ConvertibleTimestamp) time.Time {
	var ts time.Time
	switch v := tsIn.(type) {
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			log.Fatalf("Error parsing timestamp: %v", err)
		}
		ts = parsed
	case time.Time:
		ts = v
	default:
		log.Fatalf("Invalid type for timestamp: %T", v)
	}

	// Truncate to millisecond resolution.
	ts = ts.Truncate(time.Millisecond)

	if ts.Location() == time.UTC {
		// Just ensuring we log a note.
		log.Printf("timestamp without timezone found, using UTC: %v", ts)
		ts = ts.UTC()
	}
	return ts
}

// parseDuration converts a float64 or time.Duration to time.Duration.
func parseDuration(dur Duration) time.Duration {
	switch v := dur.(type) {
	case time.Duration:
		return v
	case float64:
		// interpret as "seconds"
		return time.Duration(v * float64(time.Second))
	default:
		log.Fatalf("Couldn't parse duration of invalid type %T", v)
	}
	return 0
}

// ToJSONDict merges the event's Data JSON into a map and adds "id", "timestamp", "duration".
func (e *Event) ToJSONDict() map[string]interface{} {
	// Unmarshal e.Data (datatypes.JSON) into a map.
	var dataMap map[string]interface{}
	if err := json.Unmarshal(e.Data, &dataMap); err != nil {
		// fallback if unmarshal fails
		dataMap = make(map[string]interface{})
	}

	// Create a result map with extra fields
	jsonData := make(map[string]interface{}, len(dataMap)+3)
	for k, v := range dataMap {
		jsonData[k] = v
	}
	jsonData["id"] = e.ID
	jsonData["timestamp"] = e.Timestamp.Format(time.RFC3339)
	jsonData["duration"] = e.Duration.Seconds()
	return jsonData
}

// DataEqual checks if the event's Data is equal to another map.
func (e *Event) DataEqualEvent(other *Event) bool {
	var thisMap map[string]interface{}
	var otherMap map[string]interface{}

	// Unmarshal both JSON fields into maps
	if err := json.Unmarshal(e.Data, &thisMap); err != nil {
		return false
	}
	if err := json.Unmarshal(other.Data, &otherMap); err != nil {
		return false
	}

	// Compare shallowly (keys and values)
	if len(thisMap) != len(otherMap) {
		return false
	}
	for k, v := range thisMap {
		if ov, ok := otherMap[k]; !ok || v != ov {
			return false
		}
	}
	return true
}

// DataEqualJSON compares the event's JSON data to another JSON blob
func (e *Event) DataEqualJSON(other datatypes.JSON) bool {
	// Unmarshal 'e.Data' into a map
	var eMap map[string]interface{}
	if err := json.Unmarshal(e.Data, &eMap); err != nil {
		return false
	}
	// Unmarshal 'other' into another map
	var otherMap map[string]interface{}
	if err := json.Unmarshal(other, &otherMap); err != nil {
		return false
	}
	// Compare eMap and otherMap. Basic approach:
	if len(eMap) != len(otherMap) {
		return false
	}
	for k, v := range eMap {
		if ov, ok := otherMap[k]; !ok || v != ov {
			return false
		}
	}
	return true
}
