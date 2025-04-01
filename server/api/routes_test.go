package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// testServerHandler simulates the server responses expected by the client.
func testServerHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/v1/v1/info" && r.Method == http.MethodGet: // Corrected path to match client call
		// Handle GetInfo: GET /api/v1/v1/info
		json.NewEncoder(w).Encode(map[string]interface{}{
			"hostname":    "testhost",
			"version":     "1.0",
			"server_name": "testserver",
		})

	case path == "/api/v1/v1/export" && r.Method == http.MethodGet:
		// Handle ExportAll: GET /api/v1/v1/export
		json.NewEncoder(w).Encode(map[string]interface{}{
			"buckets": map[string]interface{}{
				"bucket1": map[string]interface{}{
					"id":   "bucket1",
					"type": "test",
				},
				"bucket2": map[string]interface{}{
					"id":   "bucket2",
					"type": "test2",
				},
			},
		})
	case path == "/api/v1/v1/import" && r.Method == http.MethodPost:
		// Handle ImportAll: POST /api/v1/v1/import
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))

	case path == "/api/v1/v1/buckets/" && r.Method == http.MethodGet:
		// Handle GetBuckets: GET /api/v1/v1/buckets/
		json.NewEncoder(w).Encode(map[string]interface{}{
			"bucket1": map[string]interface{}{
				"id":   "bucket1",
				"type": "test",
			},
			"bucket2": map[string]interface{}{
				"id":   "bucket2",
				"type": "test2",
			},
		})

	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodGet:
		parts := strings.Split(path, "/")
		bucketID := parts[4]
		if len(parts) == 5 { // GET /api/v1/v1/buckets/{bucket_id}
			// Handle GetBucketMetadata
			if bucketID == "testbucket" {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":       "testbucket",
					"type":     "test",
					"client":   "testclient",
					"hostname": "testhost",
				})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if len(parts) == 6 && parts[5] == "export" { // GET /api/v1/v1/buckets/{bucket_id}/export
			// Handle ExportBucket
			if bucketID == "testbucket" {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":       "testbucket",
					"type":     "test",
					"client":   "testclient",
					"hostname": "testhost",
					"events": []map[string]interface{}{
						{"id": 1, "data": "event1"},
						{"id": 2, "data": "event2"},
					},
				})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if len(parts) == 6 && parts[5] == "events" { // GET /api/v1/v1/buckets/{bucket_id}/events
			// Handle GetEvents
			if bucketID == "testbucket" {
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"id": 1, "data": "event1"},
					{"id": 2, "data": "event2"},
				})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if len(parts) == 7 && parts[5] == "events" && parts[6] == "count" { // GET /api/v1/v1/buckets/{bucket_id}/events/count
			// Handle GetEventCount
			if bucketID == "testbucket" {
				json.NewEncoder(w).Encode(2)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if len(parts) == 7 && parts[5] == "events" { // GET /api/v1/v1/buckets/{bucket_id}/events/{event_id}
			// Handle GetEvent
			eventIDStr := parts[6]
			eventID, _ := strconv.Atoi(eventIDStr)
			if bucketID == "testbucket" && eventID == 1 {
				json.NewEncoder(w).Encode(map[string]interface{}{"id": 1, "data": "event1"})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}

	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodPost:
		parts := strings.Split(path, "/")
		bucketID := parts[4]
		if len(parts) == 5 { // POST /api/v1/v1/buckets/{bucket_id}
			// Handle CreateBucket
			if bucketID == "newbucket" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else if len(parts) == 6 && parts[5] == "events" { // POST /api/v1/v1/buckets/{bucket_id}/events
			// Handle CreateEvents
			if bucketID == "testbucket" {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"id": 3, "data": "newEvent"}) // Simulate single event insert
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodPut:
		parts := strings.Split(path, "/")
		bucketID := parts[4]
		if len(parts) == 5 { // PUT /api/v1/v1/buckets/{bucket_id}
			// Handle UpdateBucket
			if bucketID == "testbucket" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	case strings.HasPrefix(path, "/api/v1/v1/buckets/") && r.Method == http.MethodDelete:
		parts := strings.Split(path, "/")
		bucketID := parts[4]
		if len(parts) == 5 { // DELETE /api/v1/v1/buckets/{bucket_id}
			// Handle DeleteBucket
			if bucketID == "testbucket" && r.URL.Query().Get("force") == "1" {
				w.WriteHeader(http.StatusOK)
			} else if bucketID == "testbucket" && r.URL.Query().Get("force") != "1" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else if len(parts) == 7 && parts[5] == "events" { // DELETE /api/v1/v1/buckets/{bucket_id}/events/{event_id}
			// Handle DeleteEvent
			eventIDStr := parts[6]
			eventID, _ := strconv.Atoi(eventIDStr)
			if bucketID == "testbucket" && eventID == 1 {
				json.NewEncoder(w).Encode(map[string]bool{"success": true})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}

	default:
		// For any unmatched endpoint, return 404.
		w.WriteHeader(http.StatusNotFound)
	}
}

type Server struct {
	APIURL string
}

func newTestServer(apiURL string) *Server {
	return &Server{APIURL: apiURL}
}

func (c *Server) GetInfo() (map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/info") // Corrected path to match handler
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var info map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Server) GetBuckets() (map[string]map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var buckets map[string]map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&buckets)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

func (c *Server) GetBucketMetadata(bucketID string) (map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/" + bucketID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var metadata map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}
	return metadata, err
}

func (c *Server) CreateBucket(bucketID string, payload map[string]string) error {
	jsonPayload, _ := json.Marshal(payload)
	res, err := http.Post(c.APIURL+"/api/v1/v1/buckets/"+bucketID, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ErrStatusCode{StatusCode: res.StatusCode}
	}
	return nil
}

func (c *Server) UpdateBucket(bucketID string, payload map[string]interface{}) error {
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPut, c.APIURL+"/api/v1/v1/buckets/"+bucketID, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ErrStatusCode{StatusCode: res.StatusCode}
	}
	return nil
}

func (c *Server) DeleteBucket(bucketID string, force bool) error {
	req, err := http.NewRequest(http.MethodDelete, c.APIURL+"/api/v1/v1/buckets/"+bucketID, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	if force {
		q.Add("force", "1")
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrStatusCode{StatusCode: res.StatusCode}
	}
	return nil
}

func (c *Server) GetEvents(bucketID string) ([]map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/" + bucketID + "/events")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var events []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Server) CreateEvents(bucketID string, payload interface{}) (map[string]interface{}, error) {
	jsonPayload, _ := json.Marshal(payload)
	res, err := http.Post(c.APIURL+"/api/v1/v1/buckets/"+bucketID+"/events", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrStatusCode{StatusCode: res.StatusCode}
	}

	var eventResponse map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&eventResponse)
	if err != nil {
		return nil, err
	}
	return eventResponse, nil
}

func (c *Server) GetEventCount(bucketID string) (int, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/" + bucketID + "/events/count")
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var count int
	err = json.NewDecoder(res.Body).Decode(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Server) GetEvent(bucketID string, eventID int) (map[string]interface{}, error) {
	eventIDStr := strconv.Itoa(eventID)
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/" + bucketID + "/events/" + eventIDStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var event map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Server) DeleteEvent(bucketID string, eventID int) (map[string]bool, error) {
	eventIDStr := strconv.Itoa(eventID)
	req, err := http.NewRequest(http.MethodDelete, c.APIURL+"/api/v1/v1/buckets/"+bucketID+"/events/"+eventIDStr, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]bool
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Server) ExportAll() (map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/export")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var exportData map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&exportData)
	if err != nil {
		return nil, err
	}
	return exportData, nil
}

func (c *Server) ImportAll(payload map[string]interface{}) error {
	jsonPayload, _ := json.Marshal(payload)
	res, err := http.Post(c.APIURL+"/api/v1/v1/import", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ErrStatusCode{StatusCode: res.StatusCode}
	}
	return nil
}

func (c *Server) ExportBucket(bucketID string) (map[string]interface{}, error) {
	res, err := http.Get(c.APIURL + "/api/v1/v1/buckets/" + bucketID + "/export")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var exportData map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&exportData)
	if err != nil {
		return nil, err
	}
	return exportData, nil
}

type ErrStatusCode struct {
	StatusCode int
}

func (e ErrStatusCode) Error() string {
	return "status code error: " + strconv.Itoa(e.StatusCode)
}

func TestGetInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	info, err := client.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo error: %v", err)
	}
	if info["hostname"] != "testhost" {
		t.Errorf("expected hostname: testhost, got: %v", info["hostname"])
	}
	if info["version"] != "1.0" {
		t.Errorf("expected version: 1.0, got: %v", info["version"])
	}
	if info["server_name"] != "testserver" {
		t.Errorf("expected server_name: testserver, got: %v", info["server_name"])
	}
}

func TestGetBuckets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	buckets, err := client.GetBuckets()
	if err != nil {
		t.Fatalf("GetBuckets error: %v", err)
	}
	if len(buckets) != 2 {
		t.Errorf("expected 2 buckets, got: %d", len(buckets))
	}
	if buckets["bucket1"]["type"] != "test" {
		t.Errorf("expected bucket1 type: test, got: %v", buckets["bucket1"]["type"])
	}
	if buckets["bucket2"]["type"] != "test2" {
		t.Errorf("expected bucket2 type: test2, got: %v", buckets["bucket2"]["type"])
	}
}

func TestCreateBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	payload := map[string]string{"type": "test", "client": "testclient", "hostname": "testhost"}
	err := client.CreateBucket("newbucket", payload)
	if err != nil {
		t.Fatalf("CreateBucket error: %v", err)
	}
}

func TestUpdateBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	payload := map[string]interface{}{"type": "updatedTest"}
	err := client.UpdateBucket("testbucket", payload)
	if err != nil {
		t.Fatalf("UpdateBucket error: %v", err)
	}
}

func TestDeleteBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	err := client.DeleteBucket("testbucket", true)
	if err != nil {
		t.Fatalf("DeleteBucket error: %v", err)
	}
}

func TestGetEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	events, err := client.GetEvents("testbucket")
	if err != nil {
		t.Fatalf("GetEvents error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got: %d", len(events))
	}
	if events[0]["data"] != "event1" {
		t.Errorf("expected event data: event1, got: %v", events[0]["data"])
	}
}

func TestExportAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	exportData, err := client.ExportAll()
	if err != nil {
		t.Fatalf("ExportAll error: %v", err)
	}
	buckets, ok := exportData["buckets"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'buckets' key in export data")
	}
	if len(buckets) != 2 {
		t.Errorf("expected 2 exported buckets, got: %d", len(buckets))
	}
	if buckets["bucket1"].(map[string]interface{})["type"] != "test" {
		t.Errorf("expected bucket1 type: test, got: %v", buckets["bucket1"].(map[string]interface{})["type"])
	}
}

func TestImportAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer ts.Close()

	client := newTestServer(ts.URL)
	payload := map[string]interface{}{
		"buckets": map[string]interface{}{
			"bucket1": map[string]interface{}{
				"id":   "bucket1",
				"type": "test",
			},
		},
	}
	err := client.ImportAll(payload)
	if err != nil {
		t.Fatalf("ImportAll error: %v", err)
	}
}
