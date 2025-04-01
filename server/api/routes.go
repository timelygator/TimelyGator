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

func RegisterRoutes(cfg types.Config, datastore *database.Datastore, r *mux.Router) {
	api = API{
		config: &cfg,
		ds:     datastore,
	}
	r.HandleFunc("/v1/info", getInfo).Methods("GET")
	r.HandleFunc("/v1/export", export).Methods("GET")
	r.HandleFunc("/v1/import", importer).Methods("POST")

	r.HandleFunc("/v1/buckets/", getBuckets).Methods("GET")
	r.HandleFunc("/v1/buckets/{bucket_id}", bucket).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/v1/buckets/{bucket_id}/events", event).Methods("GET", "POST")
	r.HandleFunc("/v1/buckets/{bucket_id}/events/count", getCount).Methods("GET")
	r.HandleFunc("/v1/buckets/{bucket_id}/events/{event_id}", getEvent).Methods("GET", "DELETE")
	// r.HandleFunc("/v1/buckets/{bucket_id}/heartbeat", heartbeat).Methods("POST")
	r.HandleFunc("/v1/buckets/{bucket_id}/export", exportB).Methods("GET")
}

// GetInfo godoc
// @Summary Get server info
// @Description Get information about the server, such as version and build time.
// @Tags info
// @Produce json
// @Success 200 {object} types.InfoResponse
// @Failure 500 {object} types.HTTPError
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
// @Summary Get all buckets
// @Description Retrieve a list of all buckets.
// @Tags buckets
// @Produce json
// @Success 200 {object} []models.Bucket
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/ [get]
func getBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := api.GetBuckets()
	if err != nil {
		errors.HttpError(w, err, http.StatusInternalServerError)
		return
	}
	errors.JsonOK(w, buckets)
}

// GetBucketMetadata godoc
// @Summary Get bucket metadata
// @Description Retrieve metadata for a specific bucket.
// @Tags buckets
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Success 200 {object} models.Bucket
// @Failure 404 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id} [get]
// CreateBucket godoc
// @Summary Create a new bucket
// @Description Create a new bucket with the given ID and metadata.
// @Tags buckets
// @Accept json
// @Param bucket_id path string true "Bucket ID"
// @Param body body types.BucketCreationPayload true "Bucket creation payload"
// @Success 200 {string} Success
// @Failure 400 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id} [post]
// UpdateBucket godoc
// @Summary Update a bucket
// @Description Update metadata for a specific bucket.
// @Tags buckets
// @Accept json
// @Param bucket_id path string true "Bucket ID"
// @Param body body types.BucketUpdatePayload true "Bucket update payload"
// @Success 200 {string} Success
// @Failure 400 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id} [put]
// DeleteBucket godoc
// @Summary Delete a bucket
// @Description Delete a specific bucket. Requires force=1 query parameter unless testing.
// @Tags buckets
// @Param bucket_id path string true "Bucket ID"
// @Param force query string false "Force delete (required unless testing)"
// @Success 200 {string} Success
// @Failure 401 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
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

// GetEvents godoc
// @Summary Get events for a bucket
// @Description Retrieve events for a specific bucket, optionally with limit and time range.
// @Tags events
// @Produce json
// @Param bucket_id path string true "Bucket ID"
// @Param limit query int false "Limit the number of events returned"
// @Param start query string false "Start time in ISO8601 format"
// @Param end query string false "End time in ISO8601 format"
// @Success 200 {object} []models.Event
// @Failure 500 {object} types.HTTPError
// @Router /v1/buckets/{bucket_id}/events [get]
// CreateEvents godoc
// @Summary Create events for a bucket
// @Description Create one or more events in a specific bucket. Accepts single event object or array of event objects.
// @Tags events
// @Accept json
// @Param bucket_id path string true "Bucket ID"
// @Param body body object true "Event payload (single event object or array of event objects)"
// @Success 200 {object} models.Event
// @Failure 400 {object} types.HTTPError
// @Failure 500 {object} types.HTTPError
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

// // Heartbeat godoc
// // @Summary Send a heartbeat event
// // @Description Send a heartbeat event for a bucket to indicate it's still active.
// // @Tags buckets
// // @Accept json
// // @Produce json
// // @Param bucket_id path string true "Bucket ID"
// // @Param pulsetime query number true "Pulsetime in seconds"
// // @Param body body object true "Event payload"
// // @Success 200 {object} models.Event
// // @Failure 400 {object} types.HTTPError
// // @Failure 409 {object} types.HTTPError
// // @Failure 500 {object} types.HTTPError
// // @Router /v1/buckets/{bucket_id}/heartbeat [post]
// func heartbeat(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		errors.HttpErrorString(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	bucketID := mux.Vars(r)["bucket_id"]
// 	q := r.URL.Query()
// 	pulseStr := q.Get("pulsetime")
// 	if pulseStr == "" {
// 		errors.HttpErrorString(w, "Missing pulsetime", http.StatusBadRequest)
// 		return
// 	}
// 	pulsetime, err := strconv.ParseFloat(pulseStr, 64)
// 	if err != nil {
// 		errors.HttpErrorString(w, "Invalid pulsetime param", http.StatusBadRequest)
// 		return
// 	}

// 	var val map[string]interface{}
// 	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
// 		errors.HttpError(w, err, http.StatusBadRequest)
// 		return
// 	}
// 	evt, err := utils.ConvertToEvent(val)
// 	if err != nil {
// 		errors.HttpError(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	// Acquire lock
// 	locked := make(chan struct{}, 1)
// 	go func() {
// 		heartbeatLock.Lock()
// 		locked <- struct{}{}
// 	}()
// 	select {
// 	case <-locked:
// 		defer heartbeatLock.Unlock()
// 		e, err := api.Heartbeat(bucketID, evt, pulsetime)
// 		if err != nil {
// 			errors.HttpError(w, err, http.StatusInternalServerError)
// 			return
// 		}
// 		errors.JsonOK(w, e.ToJSONDict())
// 	case <-time.After(1 * time.Second):
// 		// Could not acquire lock in 1s
// 		errors.HttpErrorString(w, "Could not acquire heartbeat lock in reasonable time", http.StatusConflict)
// 	}
// }

// ExportAll godoc
// @Summary Export all buckets
// @Description Export all buckets and their data as a JSON attachment.
// @Tags export-import
// @Produce json
// @Success 200 {file} json "attachment"
// @Failure 500 {object} types.HTTPError
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
