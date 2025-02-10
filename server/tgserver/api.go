package tgserver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	db "timelygator/server/database"
	"timelygator/server/models"
	"timelygator/server/tg-core/core"
	"timelygator/server/tg-core/transform"
)

// ---------------------- Error Type for "Bucket not found" ---------------------- //
type NotFound struct {
	Code    string
	Message string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ---------------------- ServerAPI Struct ------------------------------- //
//
type ServerAPI struct {
	ds        *db.Datastore
	settings  *Settings
	testing   bool
	lastEvent map[string]*models.Event
	version   string
}

func NewServerAPI(datastore *db.Datastore, testing bool, version string) (*ServerAPI, error) {
	settings, err := NewSettings(testing)
	if err != nil {
		return nil, err
	}
	return &ServerAPI{
		ds:        datastore,
		settings:  settings,
		testing:   testing,
		lastEvent: make(map[string]*models.Event),
		version:   version,
	}, nil
}
// ---------------------- Utility Functions ------------------------------ //

//  1. Look for "device_id" file in GetDataDir("tg-server").
//  2. If it doesn’t exist, create it with a new UUID.
//  3. Return the device_id string.
func getDeviceID() (string, error) {
	dataDir := core.GetDataDir("tg-server")
	devicePath := filepath.Join(dataDir, "device_id")

	// Check if file exists
	if _, err := os.Stat(devicePath); err == nil {
		content, err := os.ReadFile(devicePath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	// Generate new UUID and write to file
	newUUID := "SOME-GENERATED-UUID" // Replace with real UUID generation
	if writeErr := os.WriteFile(devicePath, []byte(newUUID), 0600); writeErr != nil {
		return "", writeErr
	}
	return newUUID, nil
}

// checkBucketExists is a helper that checks if a bucket is known, else returns NotFound.
func (s *ServerAPI) checkBucketExists(bucketID string) error {
	bs := s.ds.Buckets() // map of ID -> metadata
	if _, ok := bs[bucketID]; !ok {
		return &NotFound{
			Code:    "NoSuchBucket",
			Message: fmt.Sprintf("No bucket named %s", bucketID),
		}
	}
	return nil
}

// parseISO8601 placeholder to parse ISO8601 strings.
func parseISO8601(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// ---------------------- ServerAPI Methods ------------------------------ //

func (s *ServerAPI) GetInfo() (map[string]interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	deviceID, err := getDeviceID()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"hostname": hostname,
		"version":  s.version,
		"testing":  s.testing,
		"device_id": deviceID,
	}, nil
}

func (s *ServerAPI) GetBuckets() (map[string]map[string]interface{}, error) {
	log.Println("Received get request for buckets")

	buckets := s.ds.Buckets() // returns map[string]map[string]interface{}
	for bID, bMeta := range buckets {
		bucket, err := s.ds.GetBucket(bID) // returns (*db.Bucket, error)
		if err != nil {
			// skip if error
			continue
		}
		lastEvents, err := bucket.Get(1, nil, nil)
		if err != nil {
			continue
		}
		if len(lastEvents) > 0 {
			lastEvent := lastEvents[0]
			endTime := lastEvent.Timestamp.Add(lastEvent.Duration)
			bMeta["last_updated"] = endTime.Format(time.RFC3339)
		}
	}
	return buckets, nil
}

func (s *ServerAPI) GetBucketMetadata(bucketID string) (map[string]interface{}, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}
	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	return bucket.Metadata(), nil
}

func (s *ServerAPI) ExportBucket(bucketID string) (map[string]interface{}, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}

	bucketMeta, err := s.GetBucketMetadata(bucketID)
	if err != nil {
		return nil, err
	}
	allEvents, err := s.GetEvents(bucketID, -1, nil, nil)
	if err != nil {
		return nil, err
	}
	bucketMeta["events"] = allEvents

	// Scrub event IDs
	if events, ok := bucketMeta["events"].([]map[string]interface{}); ok {
		for _, e := range events {
			delete(e, "id")
		}
	}
	return bucketMeta, nil
}

func (s *ServerAPI) ExportAll() (map[string]interface{}, error) {
	buckets, err := s.GetBuckets()
	if err != nil {
		return nil, err
	}
	exported := make(map[string]interface{})
	for bID := range buckets {
		bExport, err := s.ExportBucket(bID)
		if err != nil {
			log.Printf("Error exporting bucket '%s': %v\n", bID, err)
			continue
		}
		exported[bID] = bExport
	}
	return exported, nil
}

func (s *ServerAPI) ImportBucket(bucketData map[string]interface{}) error {
	bucketID, ok := bucketData["id"].(string)
	if !ok {
		return errors.New("invalid bucket data: missing 'id'")
	}
	log.Printf("Importing bucket %s", bucketID)

	bucketType, _ := bucketData["type"].(string)
	client, _ := bucketData["client"].(string)
	hostname, _ := bucketData["hostname"].(string)

	createdRaw, _ := bucketData["created"]
	var created time.Time
	switch v := createdRaw.(type) {
	case time.Time:
		created = v
	case string:
		t, err := parseISO8601(v)
		if err != nil {
			return err
		}
		created = t
	default:
		created = time.Now()
	}

	// Create the bucket in the DB
	_, err := s.CreateBucket(bucketID, bucketType, client, hostname, &created, nil)
	if err != nil {
		return err
	}

	// Insert events
	evtsRaw, ok := bucketData["events"].([]interface{})
	if !ok {
		return errors.New("invalid bucket data: events not an array")
	}
	var evts []*models.Event
	for _, e := range evtsRaw {
		switch evtMap := e.(type) {
		case map[string]interface{}:
			delete(evtMap, "id") // remove ID
			evts = append(evts, MapToEvent(evtMap))
		default:
			return errors.New("invalid event format in import")
		}
	}
	_, err = s.CreateEvents(bucketID, evts)
	return err
}

func (s *ServerAPI) ImportAll(buckets map[string]interface{}) error {
	for bucketID, bucketRaw := range buckets {
		bucketData, ok := bucketRaw.(map[string]interface{})
		if !ok {
			log.Printf("Skipping malformed bucket: %s\n", bucketID)
			continue
		}
		if err := s.ImportBucket(bucketData); err != nil {
			log.Printf("Error importing bucket '%s': %v\n", bucketID, err)
		}
	}
	return nil
}

// CreateBucket
func (s *ServerAPI) CreateBucket(
	bucketID, eventType, client, hostname string,
	created *time.Time,
	data map[string]interface{},
) (bool, error) {

	if created == nil {
		now := time.Now()
		created = &now
	}
	// If it already exists => return false
	bs := s.ds.Buckets()
	if _, found := bs[bucketID]; found {
		return false, nil
	}
	if hostname == "!local" {
		info, err := s.GetInfo()
		if err != nil {
			return false, err
		}
		if data == nil {
			data = map[string]interface{}{}
		}
		data["device_id"] = info["device_id"]
		hostname = info["hostname"].(string)
	}

	_, err := s.ds.CreateBucket(bucketID, eventType, client, hostname, *created, nil, data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *ServerAPI) UpdateBucket(
	bucketID string,
	eventType, client, hostname *string,
	data map[string]interface{},
) error {
	if err := s.checkBucketExists(bucketID); err != nil {
		return err
	}
	// Datastore's UpdateBucket expects a map of string -> interface{} for the updates
	updates := map[string]interface{}{}
	if eventType != nil {
		updates["type"] = *eventType
	}
	if client != nil {
		updates["client"] = *client
	}
	if hostname != nil {
		updates["hostname"] = *hostname
	}
	if data != nil {
		updates["datastr"] = data // or "data"
	}
	return s.ds.UpdateBucket(bucketID, updates)
}

func (s *ServerAPI) DeleteBucket(bucketID string) error {
	if err := s.checkBucketExists(bucketID); err != nil {
		return err
	}
	err := s.ds.DeleteBucket(bucketID)
	if err == nil {
		log.Printf("Deleted bucket '%s'\n", bucketID)
	}
	return err
}

func (s *ServerAPI) GetEvent(bucketID string, eventID int) (map[string]interface{}, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}
	log.Printf("Received get request for event %d in bucket '%s'\n", eventID, bucketID)

	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	event, err := bucket.GetByID(eventID)
	if err != nil || event == nil {
		return nil, err
	}
	return event.ToJSONDict(), nil
}

func (s *ServerAPI) GetEvents(bucketID string, limit int, start, end *time.Time) ([]map[string]interface{}, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}
	log.Printf("Received get request for events in bucket '%s'\n", bucketID)

	if limit == 0 {
		limit = -1
	}
	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	events, err := bucket.Get(limit, start, end)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}
	for _, e := range events {
		results = append(results, e.ToJSONDict())
	}
	return results, nil
}

func (s *ServerAPI) CreateEvents(bucketID string, events []*models.Event) (*models.Event, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}
	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	insertedEvent, err := bucket.Insert(events) // Insert(interface{})
	if err != nil {
		return nil, err
	}
	return insertedEvent, nil
}

func (s *ServerAPI) GetEventCount(bucketID string, start, end *time.Time) (int, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return 0, err
	}
	log.Printf("Received get request for eventcount in bucket '%s'\n", bucketID)

	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return 0, err
	}
	return bucket.GetEventCount(start, end)
}

func (s *ServerAPI) DeleteEvent(bucketID string, eventID int) (bool, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return false, err
	}
	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return false, err
	}
	return bucket.Delete(eventID)
}

func (s *ServerAPI) Heartbeat(bucketID string, heartbeat *models.Event, pulseTime float64) (*models.Event, error) {
	if err := s.checkBucketExists(bucketID); err != nil {
		return nil, err
	}
	log.Printf("Received heartbeat in bucket '%s'\n\ttimestamp: %v, duration: %v, pulsetime: %f\n\tdata: %+v\n",
		bucketID, heartbeat.Timestamp, heartbeat.Duration, pulseTime, heartbeat.Data)

	var lastEvent *models.Event
	lastEvent = s.lastEvent[bucketID] // from memory
	if lastEvent == nil {
		bucket, err := s.ds.GetBucket(bucketID)
		if err != nil {
			return nil, err
		}
		evts, err := bucket.Get(1, nil, nil)
		if err == nil && len(evts) > 0 {
			lastEvent = evts[0]
		}
	}

	if lastEvent != nil {
		if lastEvent.DataEqual(heartbeat.Data) {
			// Merge
			merged := transform.HeartbeatMerge(*lastEvent, *heartbeat, pulseTime)
			if merged != nil {
				log.Printf("Merging heartbeat in bucket '%s'\n", bucketID)
				s.lastEvent[bucketID] = merged
				bucket, err := s.ds.GetBucket(bucketID)
				if err != nil {
					return nil, err
				}
				if err := bucket.ReplaceLast(merged); err != nil {
					return nil, err
				}
				return merged, nil
			}
			log.Printf("Heartbeat outside pulse window, inserting new event. (bucket: %s)\n", bucketID)
		} else {
			log.Printf("Heartbeat data differs, inserting new event. (bucket: %s)\n", bucketID)
		}
	} else {
		log.Printf("Received heartbeat, but bucket was empty, inserting new event. (bucket: %s)\n", bucketID)
	}

	// Insert as new event
	bucket, err := s.ds.GetBucket(bucketID)
	if err != nil {
		return nil, err
	}
	_, insertErr := bucket.Insert([]*models.Event{heartbeat})
	if insertErr != nil {
		return nil, insertErr
	}
	s.lastEvent[bucketID] = heartbeat
	return heartbeat, nil
}

func (s *ServerAPI) GetSetting(key string) interface{} {
	return s.settings.Get(key, nil)
}

func (s *ServerAPI) SetSetting(key string, value interface{}) interface{} {
	s.settings.Set(key, value)
	return value
}
// MapToEvent is a helper to create an Event from a map[string]interface{}.
func MapToEvent(m map[string]interface{}) *models.Event {
	evt := &models.Event{}
	// Fill in evt from m as needed
	return evt
}

func (s *ServerAPI) IsTesting() bool {
    return s.testing
}
