package core

import (
	"log"
	"time"
	"gorm.io/gorm"
)

type ConvertibleTimestamp interface{}
type Duration interface{}
type Data map[string]interface{}

type Event struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp time.Time      `gorm:"not null;type:timestamp" json:"timestamp"`
	Duration  time.Duration  `gorm:"not null" json:"duration"`
	Data      Data           `gorm:"type:json" json:"data"`
}

func NewEvent(id uint, timestamp ConvertibleTimestamp, duration Duration, data Data) *Event {
	if timestamp == nil {
		log.Printf("Event initializer did not receive a timestamp argument, using now as timestamp")
		timestamp = time.Now().In(time.UTC)
	} else {
		timestamp = timestampParse(timestamp)
	}

	if data == nil {
		data = make(Data)
	}

	return &Event{
		ID:        id,
		Timestamp: timestamp.(time.Time),
		Duration:  parseDuration(duration),
		Data:      data,
	}
}

func timestampParse(tsIn ConvertibleTimestamp) time.Time {
	var ts time.Time
	switch v := tsIn.(type) {
	case string:
		var err error
		ts, err = time.Parse(time.RFC3339, v)
		if err != nil {
			log.Fatalf("Error parsing timestamp: %v", err)
		}
	case time.Time:
		ts = v
	default:
		log.Fatalf("Invalid type for timestamp: %T", v)
	}

	// Set resolution to milliseconds instead of microseconds
	ts = ts.Truncate(time.Millisecond)

	// Add timezone if not set
	if ts.Location() == time.UTC {
		log.Printf("timestamp without timezone found, using UTC: %v", ts)
		ts = ts.In(time.UTC)
	}
	return ts
}


func parseDuration(dur Duration) time.Duration {
	switch v := dur.(type) {
	case time.Duration:
		return v
	case float64:
		return time.Duration(v) * time.Second
	default:
		log.Fatalf("Couldn't parse duration of invalid type %T", v)
	}
	return 0
}
