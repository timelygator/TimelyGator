package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
	"gorm.io/datatypes"

	"timelygator/server/models"
	"timelygator/server/tg-core/core"
)

// Datastore is your high-level DB manager
type Datastore struct {
	logger      *log.Logger
	db          *gorm.DB
	testing     bool
	extraParams map[string]interface{}
}

// NewDatastore opens a SQLite DB with GORM, auto-migrates the models, and returns a Datastore.
func NewDatastore(testing bool, extraParams map[string]interface{}, logger *log.Logger) (*Datastore, error) {
	if logger == nil {
		logger = log.Default()
	}

	sqlitePath, err := buildSQLitePath(testing)
	if err != nil {
		return nil, err
	}

	// // Optionally set GORM config; e.g. log mode:
	// gormConfig := &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent), // or Info, Debug if you want logs
	// }

	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db with gorm: %w", err)
	}

	// Auto-migrate Bucket and Event
	if err := db.AutoMigrate(&models.Bucket{}, &models.Event{}); err != nil {
		return nil, fmt.Errorf("auto-migrate error: %w", err)
	}

	logger.Printf("Using GORM-based SQLite at: %s", sqlitePath)

	ds := &Datastore{
		logger:      logger,
		db:          db,
		testing:     testing,
		extraParams: extraParams,
	}
	return ds, nil
}

// buildSQLitePath decides the DB file location based on "testing" etc.
func buildSQLitePath(testing bool) (string, error) {
	dir := core.GetDataDir("tg-server")
	EnsurePathExists(dir)

	suffix := ""
	if testing {
		suffix = "-testing"
	}
	filename := fmt.Sprintf("sqlite%s.v1.db", suffix)
	return filepath.Join(dir, filename), nil
}

// EnsurePathExists ensures the directory exists (for the DB file).
func EnsurePathExists(path string) {
	if err := os.MkdirAll(path, 0o755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", path, err)
	}
}

// ----------------------------------------------------------------------
// Bucket-Related Methods
// ----------------------------------------------------------------------

// Buckets returns a map of bucket_id -> metadata
func (ds *Datastore) Buckets() map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	var buckets []models.Bucket
	if err := ds.db.Find(&buckets).Error; err != nil {
		ds.logger.Printf("Error listing buckets: %v\n", err)
		return result
	}

	for _, b := range buckets {
		result[b.ID] = map[string]interface{}{
			"id":       b.ID,
			"name":     b.Name,
			"type":     b.Type,
			"client":   b.Client,
			"hostname": b.Hostname,
			"created":  b.Created.Format(time.RFC3339),
			// Datatypes.JSON stored in Bucket.Data => might want to unmarshal it to a map
			"data": b.Data, // you can leave it as is, or parse if needed
		}
	}
	return result
}

// CreateBucket inserts a new Bucket into the DB
func (ds *Datastore) CreateBucket(
	bucketID string,
	bucketType string,
	client string,
	hostname string,
	created time.Time,
	name *string,
	data map[string]interface{},
) (*Bucket, error) {

	// Convert the incoming map[string]interface{} -> datatypes.JSON
	jsonData, err := mapToJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	newBucket := models.Bucket{
		ID:       bucketID,
		Type:     bucketType,
		Client:   client,
		Hostname: hostname,
		Created:  created,
		Name:     name,
		Data:     jsonData, // store as JSON
	}

	if err := ds.db.Create(&newBucket).Error; err != nil {
		return nil, err
	}
	return NewBucket(ds, bucketID), nil
}

// UpdateBucket updates the specified bucket in DB
func (ds *Datastore) UpdateBucket(bucketID string, updates map[string]interface{}) error {
	var existing models.Bucket
	if err := ds.db.First(&existing, "id = ?", bucketID).Error; err != nil {
		return err
	}

	// If "type" is set
	if v, ok := updates["type"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Type = vs
		}
	}
	// If "client" is set
	if v, ok := updates["client"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Client = vs
		}
	}
	// If "hostname" is set
	if v, ok := updates["hostname"]; ok {
		if vs, _ := v.(string); vs != "" {
			existing.Hostname = vs
		}
	}
	// If data is set
	if v, ok := updates["datastr"]; ok {
		if dataMap, _ := v.(map[string]interface{}); dataMap != nil {
			jsonData, err := mapToJSON(dataMap)
			if err != nil {
				return fmt.Errorf("failed to convert datastr to JSON: %w", err)
			}
			existing.Data = jsonData
		}
	}

	return ds.db.Save(&existing).Error
}

// DeleteBucket removes the bucket row from DB
func (ds *Datastore) DeleteBucket(bucketID string) error {
	return ds.db.Where("id = ?", bucketID).Delete(&models.Bucket{}).Error
}

// GetBucket checks existence and returns a handle to operate on the bucket
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

// ----------------------------------------------------------------------
// Bucket struct for runtime logic
// ----------------------------------------------------------------------
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

// Metadata loads the bucket row from DB and returns a metadata map
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
		"data":     bucket.Data, // still stored as datatypes.JSON
	}
}

// Get retrieves events for this bucket with optional limit/time filters
func (b *Bucket) Get(limit int, starttime, endtime *time.Time) ([]*models.Event, error) {
	if starttime != nil {
		ms := starttime.UnixNano() / int64(time.Millisecond)
		*starttime = time.Unix(0, ms*int64(time.Millisecond))
	}
	if endtime != nil {
		ms := endtime.UnixNano()/int64(time.Millisecond) + 1
		*endtime = time.Unix(0, ms*int64(time.Millisecond))
	}

	dbq := b.ds.db.Model(&models.Event{}).
		Where("bucket_id = ?", b.bucketID)

	if starttime != nil && !starttime.IsZero() && endtime != nil && !endtime.IsZero() {
		dbq = dbq.Where("timestamp >= ? AND timestamp <= ?", *starttime, *endtime)
	}

	if limit > 0 {
		dbq = dbq.Limit(limit)
	}

	dbq = dbq.Order("timestamp DESC")

	var events []*models.Event
	if err := dbq.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// GetByID loads a single event from DB
func (b *Bucket) GetByID(eventID int) (*models.Event, error) {
	var evt models.Event
	if err := b.ds.db.First(&evt, "id = ? AND bucket_id = ?", eventID, b.bucketID).Error; err != nil {
		return nil, err
	}
	return &evt, nil
}

// GetEventCount
func (b *Bucket) GetEventCount(starttime, endtime *time.Time) (int, error) {
	dbq := b.ds.db.Model(&models.Event{}).Where("bucket_id = ?", b.bucketID)
	if starttime != nil && !starttime.IsZero() && endtime != nil && !endtime.IsZero() {
		dbq = dbq.Where("timestamp >= ? AND timestamp <= ?", *starttime, *endtime)
	}

	var count int64
	if err := dbq.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// Insert can handle one or multiple events
func (b *Bucket) Insert(events interface{}) (*models.Event, error) {
	switch ev := events.(type) {
	case *models.Event:
		ev.BucketID = b.bucketID // ensure correct BucketID
		if err := b.ds.db.Create(ev).Error; err != nil {
			return nil, err
		}
		return ev, nil

	case []*models.Event:
		for _, e := range ev {
			e.BucketID = b.bucketID
		}
		if len(ev) == 0 {
			return nil, nil
		}
		if err := b.ds.db.Create(&ev).Error; err != nil {
			return nil, err
		}
		return nil, nil

	case []models.Event:
		for i := range ev {
			ev[i].BucketID = b.bucketID
		}
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

// Delete removes a specific event
func (b *Bucket) Delete(eventID int) (bool, error) {
	res := b.ds.db.Where("bucket_id = ?", b.bucketID).Delete(&models.Event{}, eventID)
	if res.Error != nil {
		return false, res.Error
	}
	return (res.RowsAffected == 1), nil
}

// ReplaceLast updates the most recent event in the bucket
func (b *Bucket) ReplaceLast(event *models.Event) error {
	var last models.Event
	if err := b.ds.db.Where("bucket_id = ?", b.bucketID).
		Order("timestamp desc").
		First(&last).Error; err != nil {
		return err
	}
	// Update with new data
	last.Timestamp = event.Timestamp
	last.Duration = event.Duration
	last.Data = event.Data
	return b.ds.db.Save(&last).Error
}

// Replace updates an event by ID
func (b *Bucket) Replace(eventID int, event *models.Event) error {
	var existing models.Event
	if err := b.ds.db.Where("bucket_id = ?", b.bucketID).
		First(&existing, "id = ?", eventID).Error; err != nil {
		return err
	}
	existing.Timestamp = event.Timestamp
	existing.Duration = event.Duration
	existing.Data = event.Data
	return b.ds.db.Save(&existing).Error
}

// ----------------------------------------------------------------------
// Helper: mapToJSON
// ----------------------------------------------------------------------

// mapToJSON converts a map into datatypes.JSON for storing in GORM.
func mapToJSON(data map[string]interface{}) (datatypes.JSON, error) {
	if data == nil {
		return datatypes.JSON([]byte("{}")), nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
