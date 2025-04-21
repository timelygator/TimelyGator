package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"timelygator/server/database"
	"timelygator/server/database/models"
	"timelygator/server/middleware/errors"
	"timelygator/server/utils"
	"timelygator/server/utils/types"

	"github.com/gorilla/mux"
)

var heartbeatLock sync.Mutex
var api API

// @title TimelyGator API
// @version 1.0
// @description TimelyGator is a time-tracking and activity monitoring service that provides REST APIs for managing buckets and events.
// @BasePath /v1

func RegisterRoutes(cfg types.Config, datastore *database.Datastore, r *mux.Router) {
	api = API{
		config:    &cfg,
		ds:        datastore,
		lastEvent: make(map[string]*models.Event),
	}
	r.HandleFunc("/v1/info", getInfo).Methods("GET")
	r.HandleFunc("/v1/export", export).Methods("GET")
	r.HandleFunc("/v1/import", importer).Methods("POST")

	r.HandleFunc("/v1/buckets/", getBuckets).Methods("GET")
	r.HandleFunc("/v1/buckets/{bucket_id}", bucket).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/v1/buckets/{bucket_id}/events", event).Methods("GET", "POST")
	r.HandleFunc("/v1/buckets/{bucket_id}/events/count", getCount).Methods("GET")
	r.HandleFunc("/v1/buckets/{bucket_id}/events/{event_id}", getEvent).Methods("GET", "DELETE")
	r.HandleFunc("/v1/buckets/{bucket_id}/heartbeat", heartbeat).Methods("POST")
	r.HandleFunc("/v1/buckets/{bucket_id}/export", exportB).Methods("GET")
}

// GetInfo godoc
// @Summary Get server information
// @Description Returns detailed information about the TimelyGator server instance including version,
// @Description build time, and other deployment-specific configuration.
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} types.InfoResponse "Server information retrieved successfully"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/info [get]
func getInfo(w http.ResponseWriter, r *http.Request) {
	info, err := api.GetInfo()
	if err != nil {
		errors.HttpError(w, err, http.StatusInternalServerError)
		return
	}
	errors.JsonOK(w, info)
}

// GetBuckets godoc
// @Summary List all buckets
// @Description Retrieves a list of all buckets in the system. Each bucket represents a collection
// @Description of related events and contains metadata about the tracking session.
// @Tags buckets
// @Accept json
// @Produce json
// @Success 200 {array} models.Bucket "List of buckets retrieved successfully"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/buckets [get]
func getBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := api.GetBuckets()
	if err != nil {
		errors.HttpError(w, err, http.StatusInternalServerError)
		return
	}
	errors.JsonOK(w, buckets)
}

// Bucket operations godoc
// @Summary Manage bucket operations
// @Description Endpoint for creating, retrieving, updating, and deleting buckets
// @Tags buckets
// @Accept json
// @Produce json
// @Param bucket_id path string true "Unique identifier for the bucket"
// @Param force query string false "Force deletion flag (required for DELETE unless in testing mode)"
// @Success 200 {object} models.Bucket "Operation completed successfully"
// @Success 204 {string} string "No content (for successful updates)"
// @Failure 400 {object} types.HTTPError "Invalid request parameters"
// @Failure 404 {object} types.HTTPError "Bucket not found"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/buckets/{bucket_id} [get]
// @Router /v1/buckets/{bucket_id} [post]
// @Router /v1/buckets/{bucket_id} [put]
// @Router /v1/buckets/{bucket_id} [delete]
func bucket(w http.ResponseWriter, r *http.Request) {
	bucketID := mux.Vars(r)["bucket_id"]
	switch r.Method {
	case "GET":
		// Get bucket metadata
		meta, err := api.GetBucketMetadata(bucketID)
		if err != nil {
			// handle 404 or other error
			if utils.IsNotFound(err) {
				errors.HttpError(w, err, http.StatusNotFound)
			} else {
				errors.HttpError(w, err, http.StatusInternalServerError)
			}
			return
		}
		errors.JsonOK(w, meta)

	case "POST":
		// Create bucket
		var payload struct {
			Client   string `json:"client"`
			Type     string `json:"type"`
			Hostname string `json:"hostname"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errors.HttpError(w, err, http.StatusBadRequest)
			return
		}
		created, err := api.CreateBucket(bucketID, payload.Type, payload.Client, payload.Hostname, nil, nil)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		if created {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}

	case "PUT":
		// Update bucket
		var payload struct {
			Client   *string                `json:"client"`
			Type     *string                `json:"type"`
			Hostname *string                `json:"hostname"`
			Data     map[string]interface{} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			errors.HttpError(w, err, http.StatusBadRequest)
			return
		}
		err := api.UpdateBucket(
			bucketID,
			payload.Type,
			payload.Client,
			payload.Hostname,
			payload.Data,
		)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case "DELETE":
		// Delete bucket
		q := r.URL.Query()
		force := q.Get("force")
		if api.config.Environment != "testing" {
			if force != "1" {
				errors.HttpErrorString(w, "Deleting buckets is only permitted if testing or ?force=1", http.StatusUnauthorized)
				return
			}
		}
		err := api.DeleteBucket(bucketID)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Event operations godoc
// @Summary Manage events within a bucket
// @Description Endpoint for creating and retrieving events associated with a specific bucket.
// @Description Events represent individual time-tracking entries or activities.
// @Tags events
// @Accept json
// @Produce json
// @Param bucket_id path string true "ID of the bucket containing the events"
// @Param limit query integer false "Maximum number of events to return (for GET)"
// @Param start query string false "Start time in ISO8601 format (for GET)"
// @Param end query string false "End time in ISO8601 format (for GET)"
// @Param event body object false "Event object or array of event objects (for POST)"
// @Success 200 {array} models.Event "Events retrieved/created successfully"
// @Success 201 {object} models.Event "Event created successfully"
// @Failure 400 {object} types.HTTPError "Invalid request parameters"
// @Failure 404 {object} types.HTTPError "Bucket not found"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/buckets/{bucket_id}/events [get]
// @Router /v1/buckets/{bucket_id}/events [post]
func event(w http.ResponseWriter, r *http.Request) {
	bucketID := mux.Vars(r)["bucket_id"]
	switch r.Method {
	case "GET":
		q := r.URL.Query()
		limitStr := q.Get("limit")
		startStr := q.Get("start")
		endStr := q.Get("end")

		limit := -1
		if limitStr != "" {
			if val, err := strconv.Atoi(limitStr); err == nil {
				limit = val
			}
		}
		var startTime, endTime *time.Time
		if startStr != "" {
			t, err := utils.ParseIso8601(startStr)
			if err == nil {
				startTime = &t
			}
		}
		if endStr != "" {
			t, err := utils.ParseIso8601(endStr)
			if err == nil {
				endTime = &t
			}
		}
		events, err := api.GetEvents(bucketID, limit, startTime, endTime)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		errors.JsonOK(w, events)

	case "POST":
		// create events
		var raw interface{}
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			errors.HttpError(w, err, http.StatusBadRequest)
			return
		}
		var evts []*models.Event
		switch val := raw.(type) {
		case map[string]interface{}:
			// single event
			e, err := utils.ConvertToEvent(val)
			if err != nil {
				errors.HttpError(w, err, http.StatusBadRequest)
				return
			}
			evts = []*models.Event{e}
		case []interface{}:
			// multiple events
			for _, item := range val {
				evtMap, ok := item.(map[string]interface{})
				if !ok {
					errors.HttpErrorString(w, "Invalid event data in array", http.StatusBadRequest)
					return
				}
				e, err := utils.ConvertToEvent(evtMap)
				if err != nil {
					errors.HttpError(w, err, http.StatusBadRequest)
					return
				}
				evts = append(evts, e)
			}
		default:
			errors.HttpErrorString(w, "Invalid POST data for events", http.StatusBadRequest)
			return
		}

		inserted, err := api.CreateEvents(bucketID, evts)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		if inserted != nil {
			errors.JsonOK(w, inserted.ToJSONDict())
		} else {
			w.WriteHeader(http.StatusOK) // no single event returned
		}

	default:
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GetEventCount godoc
// @Summary Get event count for a bucket
// @Description Retrieve the count of events for a specific bucket within an optional time range.
// @Tags events
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Param start query string false "Start time in ISO8601 format"
// @Param end query string false "End time in ISO8601 format"
// @Success 200 {integer} integer
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id}/events/count [get]
func getCount(w http.ResponseWriter, r *http.Request) {
	bucketID := mux.Vars(r)["bucket_id"]
	q := r.URL.Query()
	startStr := q.Get("start")
	endStr := q.Get("end")

	var startTime, endTime *time.Time
	if startStr != "" {
		t, err := utils.ParseIso8601(startStr)
		if err == nil {
			startTime = &t
		}
	}
	if endStr != "" {
		t, err := utils.ParseIso8601(endStr)
		if err == nil {
			endTime = &t
		}
	}
	count, err := api.GetEventCount(bucketID, startTime, endTime)
	if err != nil {
		errors.HttpError(w, err, http.StatusInternalServerError)
		return
	}
	errors.JsonOK(w, count)
}

// GetEvent godoc
// @Summary Get a single event
// @Description Retrieve a specific event from a bucket by its ID.
// @Tags events
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Param event_id path integer true "Event ID"
// @Success 200 {object} models.Event
// @Failure 404 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id}/events/{event_id} [get]
// DeleteEvent godoc
// @Summary Delete a single event
// @Description Delete a specific event from a bucket by its ID.
// @Tags events
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Param event_id path integer true "Event ID"
// @Success 200 {object} map[string]bool
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id}/events/{event_id} [delete]
func getEvent(w http.ResponseWriter, r *http.Request) {
	bucketID := mux.Vars(r)["bucket_id"]
	eventIDStr := mux.Vars(r)["event_id"]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		errors.HttpErrorString(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		evt, err := api.GetEvent(bucketID, eventID)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		if evt == nil {
			errors.HttpErrorString(w, "Event not found", http.StatusNotFound)
			return
		}
		errors.JsonOK(w, evt)

	case "DELETE":
		success, err := api.DeleteEvent(bucketID, eventID)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		errors.JsonOK(w, map[string]bool{"success": success})

	default:
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Heartbeat godoc
// @Summary Send bucket heartbeat
// @Description Updates or creates an event in the specified bucket to indicate active status.
// @Description If an existing event is found within the pulsetime window, it will be updated
// @Description instead of creating a new event.
// @Tags events
// @Accept json
// @Produce json
// @Param bucket_id path string true "ID of the bucket to send heartbeat to"
// @Param pulsetime query number true "Time window in seconds to merge events"
// @Param event body object true "Event data to record"
// @Success 200 {object} models.Event "Heartbeat recorded successfully"
// @Failure 400 {object} types.HTTPError "Missing or invalid parameters"
// @Failure 409 {object} types.HTTPError "Concurrent heartbeat operation in progress"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/buckets/{bucket_id}/heartbeat [post]
func heartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	bucketID := mux.Vars(r)["bucket_id"]
	q := r.URL.Query()
	pulseStr := q.Get("pulsetime")
	if pulseStr == "" {
		errors.HttpErrorString(w, "Missing pulsetime", http.StatusBadRequest)
		return
	}
	pulsetime, err := strconv.ParseFloat(pulseStr, 64)
	if err != nil {
		errors.HttpErrorString(w, "Invalid pulsetime param", http.StatusBadRequest)
		return
	}

	var val map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		errors.HttpError(w, err, http.StatusBadRequest)
		return
	}
	evt, err := utils.ConvertToEvent(val)
	if err != nil {
		errors.HttpError(w, err, http.StatusBadRequest)
		return
	}

	// Acquire lock
	locked := make(chan struct{}, 1)
	go func() {
		heartbeatLock.Lock()
		locked <- struct{}{}
	}()
	select {
	case <-locked:
		defer heartbeatLock.Unlock()
		e, err := api.Heartbeat(bucketID, evt, pulsetime)
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		errors.JsonOK(w, e.ToJSONDict())
	case <-time.After(1 * time.Second):
		// Could not acquire lock in 1s
		errors.HttpErrorString(w, "Could not acquire heartbeat lock in reasonable time", http.StatusConflict)
	}
}

// Export/Import operations godoc
// @Summary Export all bucket data
// @Description Exports all buckets and their associated events as a JSON file attachment.
// @Description The exported data can be used for backup or migration purposes.
// @Tags export-import
// @Produce json
// @Produce application/octet-stream
// @Success 200 {file} binary "JSON file containing all bucket data"
// @Failure 500 {object} types.HTTPError "Internal server error occurred"
// @Router /v1/export [get]
func export(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bucketsExport, err := api.ExportAll()
		if err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
		payload := map[string]interface{}{"buckets": bucketsExport}
		utils.WriteAttachmentJSON(w, payload, "tg-buckets-export.json")
	default:
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// ExportBucket godoc
// @Summary Export a bucket
// @Description Export a specific bucket and its data as a JSON attachment.
// @Tags export-import
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Success 200 {file} json "attachment"
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id}/export [get]
func exportB(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	bucketID := mux.Vars(r)["bucket_id"]
	bucketExport, err := api.ExportBucket(bucketID)
	if err != nil {
		errors.HttpError(w, err, http.StatusInternalServerError)
		return
	}
	idVal, ok := bucketExport["id"].(string)
	if !ok {
		// Handle the case where bucketExport["id"] isn't a string:
		// e.g., log an error, return an HTTP error, etc.
		errors.HttpErrorString(w, "Invalid bucket ID", http.StatusInternalServerError)
		return
	}

	payload := map[string]interface{}{
		"buckets": map[string]interface{}{
			idVal: bucketExport,
		},
	}
	filename := fmt.Sprintf("tg-bucket-export_%v.json", idVal)
	utils.WriteAttachmentJSON(w, payload, filename)
}

// ImportAll godoc
// @Summary Import all buckets
// @Description Import buckets and their data from a JSON payload, either as request body or multipart form.
// @Tags export-import
// @Accept json
// @Accept multipart/form-data
// @Param body body types.ImportPayload true "Import payload"
// @Success 200 {string} Success
// @Failure 400 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
// @Router /v1/import [post]
func importer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// If import comes from a form in the web-ui:
	if err := r.ParseMultipartForm(32 << 20); err == nil && len(r.MultipartForm.File) > 0 {
		for _, files := range r.MultipartForm.File {
			for _, f := range files {
				file, err := f.Open()
				if err != nil {
					errors.HttpError(w, err, http.StatusBadRequest)
					return
				}
				defer file.Close()

				var data struct {
					Buckets map[string]interface{} `json:"buckets"`
				}
				if decodeErr := json.NewDecoder(file).Decode(&data); decodeErr != nil {
					errors.HttpError(w, decodeErr, http.StatusBadRequest)
					return
				}
				if err := api.ImportAll(data.Buckets); err != nil {
					errors.HttpError(w, err, http.StatusInternalServerError)
					return
				}
			}
		}
	} else {
		// Normal import from body
		var data struct {
			Buckets map[string]interface{} `json:"buckets"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			errors.HttpError(w, err, http.StatusBadRequest)
			return
		}
		if err := api.ImportAll(data.Buckets); err != nil {
			errors.HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
