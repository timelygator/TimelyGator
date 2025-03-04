package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env"

	"timelygator/server/database/models"
	"timelygator/server/utils"
	"timelygator/server/utils/types"
)

type TimelyGatorClient struct {
	Testing         bool
	ClientName      string
	ClientHostname  string
	ServerAddress   string
	Instance        *SingleInstance
	CommitInterval  float64

	LastHeartbeat map[string]*models.Event

	// queue for requests if offline
	requestQueue *RequestQueue
}

func NewTimelyGatorClient(
	clientName string,
	testing bool,
	hostOverride *string,
	portOverride *string,
	protocol string,
) *TimelyGatorClient {
	if protocol == "" {
		protocol = "http"
	}

	var cfg types.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment config: %v", err)
	}

	serverHost := cfg.Interface
	if hostOverride != nil && *hostOverride != "" {
		serverHost = *hostOverride
	}
	serverPort := cfg.Port
	if portOverride != nil && *portOverride != "" {
		serverPort = *portOverride
	}

	serverAddress := fmt.Sprintf("%s://%s:%s", protocol, serverHost, serverPort)

	h, err := os.Hostname()
	if err != nil || h == "" {
		h = "unknown-host"
	}

	// singleinstance
	// TODO: error handling
	inst, err := NewSingleInstance(
		fmt.Sprintf("%s-at-%s-on-%s", clientName, serverHost, serverPort),
	)

	c := &TimelyGatorClient{
		Testing:        testing,
		ClientName:     clientName,
		ClientHostname: h,
		ServerAddress:  serverAddress,
		Instance:       inst,
		CommitInterval: 60.0,
		LastHeartbeat:  make(map[string]*models.Event),
	}
	c.requestQueue = NewRequestQueue(c)
	return c
}

func (c *TimelyGatorClient) _url(endpoint string) string {
	return fmt.Sprintf("%s/api/v1/v1/%s", c.ServerAddress, endpoint)
}

// GET request to fetch data from the server.
func (c *TimelyGatorClient) get(endpoint string, params map[string]string) (*http.Response, error) {
	url := c._url(endpoint)
	if len(params) > 0 {
		url = appendQuery(url, params)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("GET %s => status %d", url, resp.StatusCode)
	}
	return resp, nil
}

// POST request to send data to the server.
func (c *TimelyGatorClient) post(
	endpoint string,
	data interface{},
	params map[string]string,
) (*http.Response, error) {
	url := c._url(endpoint)
	if len(params) > 0 {
		url = appendQuery(url, params)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("POST %s => status %d", url, resp.StatusCode)
	}
	return resp, nil
}

// DELETE request to remove data from the server.
func (c *TimelyGatorClient) deleteReq(endpoint string, data interface{}) (*http.Response, error) {
	url := c._url(endpoint)
	var body []byte
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = b
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("DELETE %s => status %d", url, resp.StatusCode)
	}
	return resp, nil
}

func appendQuery(base string, params map[string]string) string {
	base += "?"
	for k, v := range params {
		base += fmt.Sprintf("%s=%s&", k, v)
	}
	return base[:len(base)-1]
}

// GetInfo fetches server info.
func (c *TimelyGatorClient) GetInfo() (map[string]interface{}, error) {
	resp, err := c.get("info", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *TimelyGatorClient) GetEvent(bucketID string, eventID int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("buckets/%s/events/%d", bucketID, eventID)
	resp, err := c.get(endpoint, nil)
	if err != nil {
		// Check if error is 404
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) GetEvents(
	bucketID string,
	limit int,
	start, end *time.Time,
) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("buckets/%s/events", bucketID)
	params := make(map[string]string)
	if limit >= 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if start != nil {
		params["start"] = start.Format(time.RFC3339)
	}
	if end != nil {
		params["end"] = end.Format(time.RFC3339)
	}
	resp, err := c.get(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) InsertEvent(bucketID string, evt interface{}) error {
	endpoint := fmt.Sprintf("buckets/%s/events", bucketID)
	data := []interface{}{evt} // single event
	_, err := c.post(endpoint, data, nil)
	return err
}

func (c *TimelyGatorClient) InsertEvents(bucketID string, evts []interface{}) error {
	endpoint := fmt.Sprintf("buckets/%s/events", bucketID)
	_, err := c.post(endpoint, evts, nil)
	return err
}

func (c *TimelyGatorClient) DeleteEvent(bucketID string, eventID int) error {
	endpoint := fmt.Sprintf("buckets/%s/events/%d", bucketID, eventID)
	_, err := c.deleteReq(endpoint, nil)
	return err
}

func (c *TimelyGatorClient) GetEventCount(
	bucketID string,
	start, end *time.Time,
) (int, error) {
	endpoint := fmt.Sprintf("buckets/%s/events/count", bucketID)
	params := make(map[string]string)
	if start != nil {
		params["start"] = start.Format(time.RFC3339)
	}
	if end != nil {
		params["end"] = end.Format(time.RFC3339)
	}
	resp, err := c.get(endpoint, params)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return 0, err
	}
	countStr := buf.String()
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// If queued=true, we do merging logic, else direct post
func (c *TimelyGatorClient) Heartbeat(
	bucketID string,
	event interface{},
	pulseTime float64,
	queued bool,
	commitInterval *float64,
) error {
	endpoint := fmt.Sprintf("buckets/%s/heartbeat?pulsetime=%f", bucketID, pulseTime)
	ci := c.CommitInterval
	if commitInterval != nil {
		ci = *commitInterval
	}

	if queued {
		// Pre-merge in memory
		last, ok := c.LastHeartbeat[bucketID]
		var newEvent *models.Event
		// Convert event interface{} => *models.Event if needed
		// or just store raw?
		ev, ok2 := event.(*models.Event)
		if !ok2 {
			log.Println("Heartbeat event not *models.Event, skipping merge.")
			_, err := c.post(endpoint, event, nil)
			return err
		}
		newEvent = ev

		if !ok {
			c.LastHeartbeat[bucketID] = newEvent
			return nil
		}
		merged := utils.HeartbeatMerge(*last, *newEvent, pulseTime)
		if merged != nil {
			diff := merged.Duration.Seconds()
			if diff >= ci {
				data := merged.ToJSONDict()
				c.requestQueue.AddRequest(endpoint, data)
				c.LastHeartbeat[bucketID] = newEvent
			} else {
				c.LastHeartbeat[bucketID] = merged
			}
		} else {
			data := last.ToJSONDict()
			c.requestQueue.AddRequest(endpoint, data)
			c.LastHeartbeat[bucketID] = newEvent
		}
		return nil
	}

	// direct post
	_, err := c.post(endpoint, event, nil)
	return err
}

// GetBucketsMap fetches all bucket metadata.
func (c *TimelyGatorClient) GetBucketsMap() (map[string]interface{}, error) {
	resp, err := c.get("buckets/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) CreateBucket(bucketID, eventType string, queued bool) error {
	if queued {
		c.requestQueue.RegisterBucket(bucketID, eventType)
		return nil
	}

	endpoint := fmt.Sprintf("buckets/%s", bucketID)
	data := map[string]interface{}{
		"client":   c.ClientName,
		"hostname": c.ClientHostname,
		"type":     eventType,
	}
	_, err := c.post(endpoint, data, nil)
	log.Printf("BUCKET CREATED")
	return err
}

func (c *TimelyGatorClient) DeleteBucket(bucketID string, force bool) error {
	endpoint := fmt.Sprintf("buckets/%s", bucketID)
	if force {
		endpoint += "?force=1"
	}
	_, err := c.deleteReq(endpoint, nil)
	return err
}

func (c *TimelyGatorClient) ExportAll() (map[string]interface{}, error) {
	resp, err := c.get("export", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) ExportBucket(bucketID string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("buckets/%s/export", bucketID)
	resp, err := c.get(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) ImportBucket(bucket map[string]interface{}) error {
	endpoint := "import"

	bucketID, ok := bucket["id"].(string)
	if !ok {
		return fmt.Errorf("bucket ID is not a string: %v", bucket["id"])
	}

	data := map[string]interface{}{
		"buckets": map[string]interface{}{
			bucketID: bucket,
		},
	}

	_, err := c.post(endpoint, data, nil)
	return err
}


func (c *TimelyGatorClient) Query(
	queryStr string,
	timeperiods [][2]time.Time,
	name *string,
	cache bool,
) ([]interface{}, error) {
	endpoint := "query/"
	params := make(map[string]string)
	if cache {
		if name == nil || *name == "" {
			return nil, errors.New("not allowed to do caching without a query name")
		}
		params["name"] = *name
		params["cache"] = "1"
	}

	var tps []string
	for _, tp := range timeperiods {
		start, end := tp[0], tp[1]
		// e.g. "2023-03-20T05:00:00Z/2023-03-20T10:00:00Z"
		tps = append(tps, start.Format(time.RFC3339)+"/"+end.Format(time.RFC3339))
	}

	data := map[string]interface{}{
		"timeperiods": tps,
		"query":       strings.Split(queryStr, "\n"),
	}
	resp, err := c.post(endpoint, data, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) GetSetting(key *string) (map[string]interface{}, error) {
	endpoint := "settings"
	if key != nil && *key != "" {
		endpoint = fmt.Sprintf("settings/%s", *key)
	}
	resp, err := c.get(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *TimelyGatorClient) SetSetting(key string, value string) error {
	endpoint := fmt.Sprintf("settings/%s", key)
	_, err := c.post(endpoint, value, nil)
	return err
}

func (c *TimelyGatorClient) Connect() {
	if !c.requestQueue.IsAlive() {
		c.requestQueue.Start()
	}
}

func (c *TimelyGatorClient) Disconnect() {
	c.requestQueue.Stop()
	c.requestQueue.Wait()
	// discard old thread object, create new
	c.requestQueue = NewRequestQueue(c)
}

func (c *TimelyGatorClient) WaitForStart(timeout int) error {
    if timeout == 0 {
        timeout = 10
    }

    start := time.Now()
    sleepTime := 100 * time.Millisecond
    maxSleepTime := 2 * time.Second // Cap max sleep time to 2 seconds

    for time.Since(start).Seconds() < float64(timeout) {
        if _, err := c.GetInfo(); err == nil {
            log.Printf("[WaitForStart] Server at %s is ready.", c.ServerAddress)
            return nil
        } else {
            log.Printf("[WaitForStart] Server not ready: %v (retrying in %s)", err, sleepTime)
        }

        time.Sleep(sleepTime)
        if sleepTime < maxSleepTime {
            sleepTime *= 2
        }
    }

    return fmt.Errorf("[WaitForStart] Server at %s did not start in time", c.ServerAddress)
}


type QueuedRequest struct {
	Endpoint string
	Data     map[string]interface{}
}

type BucketReg struct {
	ID   string
	Type string
}

type RequestQueue struct {
	client *TimelyGatorClient

	connected bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.Mutex // for isAlive

	registeredBuckets []BucketReg
	attemptReconnect  time.Duration

	queueMu   sync.Mutex
	queue     []QueuedRequest
	ticker    *time.Ticker
}

func NewRequestQueue(client *TimelyGatorClient) *RequestQueue {
	return &RequestQueue{
		client:           client,
		stopChan:         make(chan struct{}),
		registeredBuckets: []BucketReg{},
		attemptReconnect:  10 * time.Second,
		queue:            []QueuedRequest{},
	}
}

func (rq *RequestQueue) IsAlive() bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	// if stopChan is closed or we never started?
	return true // if we want to check a "stopped" bool
}

func (rq *RequestQueue) Start() {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	// start background goroutine
	rq.wg.Add(1)
	rq.ticker = time.NewTicker(200 * time.Millisecond)
	go rq.run()
}

func (rq *RequestQueue) Stop() {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	close(rq.stopChan)
}

func (rq *RequestQueue) Wait() {
	rq.wg.Wait()
}

func (rq *RequestQueue) run() {
	defer rq.wg.Done()
	for {
		select {
		case <-rq.stopChan:
			log.Println("RequestQueue stopping.")
			rq.ticker.Stop()
			return
		case <-rq.ticker.C:
			if !rq.connected {
				rq.tryConnect()
			} else {
				rq.dispatch()
			}
		}
	}
}

func (rq *RequestQueue) tryConnect() {
	if err := rq.createBuckets(); err != nil {
		log.Printf("Not connected, will retry. Err=%v", err)
		rq.connected = false
	} else {
		rq.connected = true
		log.Printf("Connection established. queue size=%d", len(rq.queue))
	}
}

func (rq *RequestQueue) createBuckets() error {
	for _, b := range rq.registeredBuckets {
		if err := rq.client.CreateBucket(b.ID, b.Type, false); err != nil {
			return err
		}
	}
	return nil
}

func (rq *RequestQueue) dispatch() {
	rq.queueMu.Lock()
	defer rq.queueMu.Unlock()
	if len(rq.queue) == 0 {
		return
	}
	item := rq.queue[0]
	if err := rq.send(item); err != nil {
		rq.connected = false
		log.Printf("Failed to dispatch => %v", err)
		time.Sleep(500 * time.Millisecond)
		return
	}
	rq.queue = rq.queue[1:]
}

func (rq *RequestQueue) send(item QueuedRequest) error {
	_, err := rq.client.post(item.Endpoint, item.Data, nil)
	return err
}

func (rq *RequestQueue) addRequest(endpoint string, data map[string]interface{}) {
	rq.queueMu.Lock()
	defer rq.queueMu.Unlock()
	rq.queue = append(rq.queue, QueuedRequest{Endpoint: endpoint, Data: data})
}

func (rq *RequestQueue) AddRequest(endpoint string, data map[string]interface{}) {
	if endpoint == "" {
		log.Println("AddRequest: endpoint is empty, ignoring.")
		return
	}
	rq.addRequest(endpoint, data)
}

func (rq *RequestQueue) RegisterBucket(bucketID, eventType string) {
	rq.registeredBuckets = append(rq.registeredBuckets, BucketReg{ID: bucketID, Type: eventType})
}