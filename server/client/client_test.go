package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// testServerHandler simulates the server responses expected by the client.
func testServerHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	// Handle GetInfo: GET /api/v1/v1/info
	case path == "/api/v1/v1/info" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})

	// Handle GetEventCount: GET /api/v1/v1/buckets/{bucketID}/events/count
	case strings.Contains(path, "/events/count") && r.Method == http.MethodGet:
		io.WriteString(w, "42")

	// Handle GetEvent: GET /api/v1/v1/buckets/{bucketID}/events/{eventID}
	case strings.Contains(path, "/events/") && r.Method == http.MethodGet:
		parts := strings.Split(path, "/")
		// The event ID should be the last segment.
		eventIDStr := parts[len(parts)-1]
		id, err := strconv.Atoi(eventIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "data": "sample event"})

	// Handle GetEvents: GET /api/v1/v1/buckets/{bucketID}/events
	case strings.Contains(path, "/events") && r.Method == http.MethodGet:
		// Return a slice with sample events.
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "data": "sample1"},
			{"id": 2, "data": "sample2"},
		})

	// Handle InsertEvent and InsertEvents: POST /api/v1/v1/buckets/{bucketID}/events
	case strings.Contains(path, "/events") && r.Method == http.MethodPost:
		// Simply return a confirmation response.
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "inserted"})

	// Handle Heartbeat: POST /api/v1/v1/buckets/{bucketID}/heartbeat?pulsetime=...
	case strings.Contains(path, "/heartbeat") && r.Method == http.MethodPost:
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "heartbeat received"})

	// Handle CreateBucket: POST /api/v1/v1/buckets/{bucketID}
	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodPost && !strings.Contains(path, "events"):
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "bucket created"})

	// Handle DeleteBucket and DeleteEvent: DELETE /api/v1/v1/buckets/{bucketID}...
	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodDelete:
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "deleted"})

	// Handle ExportBucket: GET /api/v1/v1/buckets/{bucketID}/export
	case strings.Contains(path, "/export") && strings.Contains(path, "buckets/") && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(map[string]interface{}{"export": "bucket data"})

	// Handle ExportAll: GET /api/v1/v1/export
	case path == "/api/v1/v1/export" && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(map[string]interface{}{"export": "all data"})

	// Handle ImportBucket: POST /api/v1/v1/import
	case path == "/api/v1/v1/import" && r.Method == http.MethodPost:
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "imported"})

	// Handle Query: POST /api/v1/v1/query/
	case strings.HasSuffix(path, "/query/") && r.Method == http.MethodPost:
		json.NewEncoder(w).Encode([]interface{}{"result1", "result2"})

	// Handle GetSetting: GET /api/v1/v1/settings or /api/v1/v1/settings/{key}
	case strings.HasPrefix(path, "/api/v1/v1/settings") && r.Method == http.MethodGet:
		json.NewEncoder(w).Encode(map[string]interface{}{"key": "value"})

	// Handle SetSetting: POST /api/v1/v1/settings/{key}
	case strings.HasPrefix(path, "/api/v1/v1/settings") && r.Method == http.MethodPost:
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "setting updated"})

	default:
		// For any unmatched endpoint, return 404.
		w.WriteHeader(http.StatusNotFound)
	}
}

// newTestClient creates a new TimelyGatorClient whose ServerAddress is set to the test server's URL.
func newTestClient(tsURL string) *TimelyGatorClient {
	// Create the client with test parameters.
	c := NewTimelyGatorClient("test-client", true, nil, nil, "http")
	// Override the server address to point to our test server.
	c.ServerAddress = tsURL
	return c
}

func TestGetInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	info, err := client.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo error: %v", err)
	}
	if status, ok := info["status"]; !ok || status != "ok" {
		t.Errorf("Expected status 'ok', got %v", info["status"])
	}
}

func TestGetEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	event, err := client.GetEvent("bucket1", 123)
	if err != nil {
		t.Fatalf("GetEvent error: %v", err)
	}
	if id, ok := event["id"]; !ok || int(id.(float64)) != 123 {
		t.Errorf("Expected event id 123, got %v", event["id"])
	}
}

func TestGetEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	events, err := client.GetEvents("bucket1", 10, nil, nil)
	if err != nil {
		t.Fatalf("GetEvents error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestInsertEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	err := client.InsertEvent("bucket1", map[string]interface{}{"data": "event1"})
	if err != nil {
		t.Fatalf("InsertEvent error: %v", err)
	}
}

func TestInsertEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	events := []interface{}{
		map[string]interface{}{"data": "event1"},
		map[string]interface{}{"data": "event2"},
	}
	err := client.InsertEvents("bucket1", events)
	if err != nil {
		t.Fatalf("InsertEvents error: %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	err := client.DeleteEvent("bucket1", 456)
	if err != nil {
		t.Fatalf("DeleteEvent error: %v", err)
	}
}

func TestGetEventCount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	count, err := client.GetEventCount("bucket1", nil, nil)
	if err != nil {
		t.Fatalf("GetEventCount error: %v", err)
	}
	if count != 42 {
		t.Errorf("Expected count 42, got %d", count)
	}
}

func TestHeartbeat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	// Test heartbeat with queued false (direct post)
	err := client.Heartbeat("bucket1", map[string]interface{}{"data": "heartbeat event"}, 1.0, false, nil)
	if err != nil {
		t.Fatalf("Heartbeat error: %v", err)
	}
}

func TestGetBucketsMap(t *testing.T) {
	// For GetBucketsMap, our test handler for "buckets/" is not specifically defined.
	// We simulate it by returning a simple JSON object.
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/v1/buckets/" && r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]interface{}{"bucket1": "data"})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	buckets, err := client.GetBucketsMap()
	if err != nil {
		t.Fatalf("GetBucketsMap error: %v", err)
	}
	if _, ok := buckets["bucket1"]; !ok {
		t.Errorf("Expected bucket1 key in buckets map")
	}
}

func TestCreateAndDeleteBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)

	// Test CreateBucket with queued == false (direct POST)
	err := client.CreateBucket("bucket1", "eventType", false)
	if err != nil {
		t.Fatalf("CreateBucket error: %v", err)
	}

	// Test DeleteBucket with force = true
	err = client.DeleteBucket("bucket1", true)
	if err != nil {
		t.Fatalf("DeleteBucket error: %v", err)
	}
}

func TestExportAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	exportData, err := client.ExportAll()
	if err != nil {
		t.Fatalf("ExportAll error: %v", err)
	}
	if val, ok := exportData["export"]; !ok || val != "all data" {
		t.Errorf("Expected export 'all data', got %v", exportData["export"])
	}
}

func TestExportBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	exportData, err := client.ExportBucket("bucket1")
	if err != nil {
		t.Fatalf("ExportBucket error: %v", err)
	}
	if val, ok := exportData["export"]; !ok || val != "bucket data" {
		t.Errorf("Expected export 'bucket data', got %v", exportData["export"])
	}
}

func TestImportBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	// Create a dummy bucket map with an "id" field.
	bucket := map[string]interface{}{
		"id":   "bucket1",
		"data": "dummy",
	}
	err := client.ImportBucket(bucket)
	if err != nil {
		t.Fatalf("ImportBucket error: %v", err)
	}
}

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	// Create a sample time period
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	timeperiods := [][2]time.Time{{start, end}}

	queryName := "sample-query"
	results, err := client.Query("select * from events", timeperiods, &queryName, true)
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 query results, got %d", len(results))
	}
}

func TestGetAndSetSetting(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)

	// Test GetSetting without key.
	setting, err := client.GetSetting(nil)
	if err != nil {
		t.Fatalf("GetSetting error: %v", err)
	}
	if val, ok := setting["key"]; !ok || val != "value" {
		t.Errorf("Expected setting 'value', got %v", setting["key"])
	}

	// Test SetSetting.
	err = client.SetSetting("testKey", "testValue")
	if err != nil {
		t.Fatalf("SetSetting error: %v", err)
	}
}

func TestWaitForStart(t *testing.T) {
	// For WaitForStart, we simulate a server that is already ready.
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestClient(ts.URL)
	// The default timeout is 10 seconds; this test should return immediately.
	err := client.WaitForStart(5)
	if err != nil {
		t.Fatalf("WaitForStart error: %v", err)
	}
}
