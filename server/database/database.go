package db

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"timelygator/server/models"
	"timelygator/server/tg-core/core"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Datastore struct {
	logger    *log.Logger
	db        *gorm.DB
	testing   bool
	extraParams map[string]interface{}
}

// NewDatastore opens a SQLite DB with GORM, auto-migrates the models, and returns a Datastore.
func NewDatastore(testing bool, extraParams map[string]interface{}, logger *log.Logger) (*Datastore, error) {
	if logger == nil {
		logger = log.Default()
	}

	// Decides the DB filepath
	filepath, err := buildSQLitePath(testing)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db with gorm: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Bucket{}, &models.Event{}); err != nil {
		return nil, fmt.Errorf("auto-migrate error: %w", err)
	}

	logger.Printf("Using GORM-based SQLite at: %s", filepath)

	ds := &Datastore{
		logger:      logger,
		db:          db,
		testing:     testing,
		extraParams: extraParams,
	}
	return ds, nil
}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Bucket{}, &models.Event{}); err != nil {
		return nil, fmt.Errorf("auto-migrate error: %w", err)
	}

	logger.Printf("Using GORM-based SQLite at: %s", filepath)

	ds := &Datastore{
		logger:      logger,
		db:          db,
		testing:     testing,
		extraParams: extraParams,
	}
	return ds, nil
}

func buildSQLitePath(testing bool) (string, error) {
	dir := core.GetDataDir("tg-server")
	suffix := ""
	if testing {
		suffix = "-testing"
	}
	filename := fmt.Sprintf("sqlite%s.v1.db", suffix)
	return filepath.Join(dir, filename), nil
}

// Batches returns a map of bucket_id -> metadata. 
func (ds *Datastore) Buckets() map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	var buckets []models.Bucket
	if err := ds.db.Find(&buckets).Error; err != nil {
		ds.logger.Printf("Error listing buckets: %v\n", err)
		return result
	}

	for _, b := range buckets {
		// Convert GORM model -> metadata map
		meta := map[string]interface{}{
			"id":       b.ID,
			"name":     b.Name,
			"type":     b.Type,
			"client":   b.Client,
			"hostname": b.Hostname,
			"created":  b.Created.Format(time.RFC3339),
			"data":     b.Data,
		}
		result[b.ID] = meta
	}
	return result
}

func (ds *Datastore) CreateBucket(
	bucketID string,
	bucketType string,
	client string,
	hostname string,
	created time.Time,
	name *string,
	data map[string]interface{},
) (*Bucket, error) {
	// Insert into DB via GORM
	newBucket := models.Bucket{
		ID:        bucketID,
		Type:      bucketType,
		Client:    client,
		Hostname:  hostname,
		Created:   created,
		Name:      name,
		Data:      data,
	}

	if err := ds.db.Create(&newBucket).Error; err != nil {
		return nil, err
	}
	// Return a *Bucket object that references the GORM model.
	return NewBucket(ds, bucketID), nil
}

func (ds *Datastore) UpdateBucket(bucketID string, updates map[string]interface{}) error {
	// We find the existing row, then apply updates
	var existing models.Bucket
	if err := ds.db.First(&existing, "id = ?", bucketID).Error; err != nil {
		return err
	}

	if v, ok := updates["type"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Type = vs
		}
	}
	if v, ok := updates["client"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Client = vs
		}
	}
	if v, ok := updates["hostname"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Hostname = vs
		}
	}
	if v, ok := updates["datastr"]; ok {
		if dataMap, _ := v.(map[string]interface{}); dataMap != nil {
			existing.Data = dataMap
		}
	}

	return ds.db.Save(&existing).Error
}

// DeleteBucket
func (ds *Datastore) DeleteBucket(bucketID string) error {
	// If we want to also delete associated events, we can define a foreign-key relationship
	// with "OnDelete:CASCADE" or do a separate .Where("bucket_id=?").Delete(&Event{})
	return ds.db.Where("id = ?", bucketID).Delete(&models.Bucket{}).Error
}

// GetBucket returns the "bucket" if it exists
func (ds *Datastore) GetBucket(bucketID string) (*Bucket, error) {
	var count int64
	if err := ds.db.Model(&models.Bucket{}).
		Where("id = ?", bucketID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("bucket %q does not exist", bucketID)
	}
	return NewBucket(ds, bucketID), nil
}

// Bucket is the GORM-backed "bucket handle"
type Bucket struct {
	ds       *Datastore
	bucketID string
}

func NewBucket(ds *Datastore, bucketID string) *Bucket {
	return &Bucket{
		ds:       ds,
		bucketID: bucketID,
	}
}

// Metadata can read the bucket row from DB
func (b *Bucket) Metadata() map[string]interface{} {
	var bucket models.Bucket
	if err := b.ds.db.First(&bucket, "id = ?", b.bucketID).Error; err != nil {
		b.ds.logger.Printf("Error in Metadata() for bucket %s: %v", b.bucketID, err)
		return nil
	}
	return map[string]interface{}{
		"id":       bucket.ID,
		"name":     bucket.Name,
		"type":     bucket.Type,
		"client":   bucket.Client,
		"hostname": bucket.Hostname,
		"created":  bucket.Created.Format(time.RFC3339),
		"data":     bucket.Data,
	}
}

func (b *Bucket) Get(limit int, starttime, endtime *time.Time) ([]*models.Event, error) {
	// Round start/end times to nearest millisecond
	if starttime != nil {
		ms := starttime.UnixNano() / int64(time.Millisecond)
		*starttime = time.Unix(0, ms*int64(time.Millisecond))
	}
	if endtime != nil {
		ms := endtime.UnixNano()/int64(time.Millisecond) + 1
		*endtime = time.Unix(0, ms*int64(time.Millisecond))
	}

	// Build the base query filtering on bucket_id
	dbq := b.ds.db.Model(&models.Event{}).
		Where("bucket_id = ?", b.bucketID)

	// If start/end time are provided, filter by them
	if starttime != nil && !starttime.IsZero() && endtime != nil && !endtime.IsZero() {
		dbq = dbq.Where("timestamp >= ? AND timestamp <= ?", *starttime, *endtime)
	}

	// If limit > 0, apply it. If limit == -1, do no limit
	if limit > 0 {
		dbq = dbq.Limit(limit)
	}

	// Sort by timestamp descending
	dbq = dbq.Order("timestamp DESC")

	var events []*models.Event
	if err := dbq.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// GetByID
func (b *Bucket) GetByID(eventID int) (*models.Event, error) {
	var evt models.Event
	if err := b.ds.db.First(&evt, "id = ?", eventID).Error; err != nil {
		return nil, err
	}
	return &evt, nil
}

// GetEventCount
func (b *Bucket) GetEventCount(starttime, endtime *time.Time) (int, error) {
	dbq := b.ds.db.Model(&models.Event{})
	// Add filters if desired
	var count int64
	if err := dbq.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// Insert can handle single or multiple events
func (b *Bucket) Insert(events interface{}) (*models.Event, error) {
	switch ev := events.(type) {
	case *models.Event:
		if err := b.ds.db.Create(ev).Error; err != nil {
			return nil, err
		}
		return ev, nil
	case []*models.Event:
		if len(ev) == 0 {
			return nil, nil
		}
		if err := b.ds.db.Create(&ev).Error; err != nil {
			return nil, err
		}
		return nil, nil
	case []models.Event:
		if len(ev) == 0 {
			return nil, nil
		}
		if err := b.ds.db.Create(&ev).Error; err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid events type in Insert(...)")
	}
}

// Delete
func (b *Bucket) Delete(eventID int) (bool, error) {
	if err := b.ds.db.Delete(&models.Event{}, eventID).Error; err != nil {
		return false, err
	}
	return true, nil
}

// ReplaceLast - do a GORM query for the last event, then update it
func (b *Bucket) ReplaceLast(event *models.Event) error {
	var last models.Event
	// find the last event for this bucket
	if err := b.ds.db.Order("timestamp desc").First(&last).Error; err != nil {
		return err
	}
	// Update last with data from the new event
	last.Timestamp = event.Timestamp
	last.Duration = event.Duration
	last.Data = event.Data
	return b.ds.db.Save(&last).Error
}

// Replace replaces the event with eventID
func (b *Bucket) Replace(eventID int, event *models.Event) error {
	var existing models.Event
	if err := b.ds.db.First(&existing, "id = ?", eventID).Error; err != nil {
		return err
	}
	existing.Timestamp = event.Timestamp
	existing.Duration = event.Duration
	existing.Data = event.Data
	return b.ds.db.Save(&existing).Error
}
