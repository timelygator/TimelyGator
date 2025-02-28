package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"os"

	"github.com/caarlos0/env"

	"timelygator/server/database/models"
	"timelygator/server/utils"
	"timelygator/server/utils/types"
)

type RequestQueueItem struct {
	Endpoint string
	Data     map[string]interface{}
}

type BucketRegistration struct {
	ID   string
	Type string
}

type ActivityWatchClient struct {
	Testing        bool
	ClientName     string
	ClientHostname string
	ServerAddress  string

	Instance *SingleInstance

	CommitInterval float64

	logger        *log.Logger
	lastHeartbeat map[string]*models.Event

	requestQueue *RequestQueue
}

func NewActivityWatchClient(clientName string, testing bool) *ActivityWatchClient {
	var cfg types.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	serverHost := cfg.Interface
	serverPort := cfg.Port
	serverAddress := fmt.Sprintf("http://%s:%s", serverHost, serverPort)

	clientHostname, err := osHostname()
	if err != nil {
		clientHostname = "unknown-host"
	}

	inst, err := NewSingleInstance(
		fmt.Sprintf("%s-at-%s-on-%s", clientName, serverHost, serverPort),
	)

	c := &ActivityWatchClient{
		Testing:        testing,
		ClientName:     clientName,
		ClientHostname: clientHostname,
		ServerAddress:  serverAddress,
		Instance:       inst,
		CommitInterval: 60.0,
		logger:         log.Default(),
		lastHeartbeat:  make(map[string]*models.Event),
	}

	c.requestQueue = NewRequestQueue(c)
	return c
}

func osHostname() (string, error) {
	hn, err := os.Hostname()
	if err != nil || hn == "" {
		return "", err
	}
	return hn, nil
}

func (c *ActivityWatchClient) _url(endpoint string) string {
	return fmt.Sprintf("%s/api/0/%s", c.ServerAddress, endpoint)
}

// get does a simple GET call
func (c *ActivityWatchClient) get(endpoint string, params map[string]string) (*http.Response, error) {
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

// post does a POST call
func (c *ActivityWatchClient) post(endpoint string, data interface{}, params map[string]string) (*http.Response, error) {
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

// deleteReq does a DELETE call
func (c *ActivityWatchClient) deleteReq(endpoint string, data interface{}) (*http.Response, error) {
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

// appendQuery is a helper to attach params to a URL
func appendQuery(base string, params map[string]string) string {
	query := "?"
	for k, v := range params {
		query += fmt.Sprintf("%s=%s&", k, v)
	}
	return base + query[:len(query)-1]
}

// GetInfo => GET /info
func (c *ActivityWatchClient) GetInfo() (map[string]interface{}, error) {
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

// GetBuckets => GET /buckets/
func (c *ActivityWatchClient) GetBuckets() (map[string]interface{}, error) {
	resp, err := c.get("buckets/", nil)
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

// CreateBucket => POST /buckets/{bucket_id}
func (c *ActivityWatchClient) CreateBucket(bucketID, eventType string, queued bool) error {
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
	return err
}

// DeleteBucket => DELETE /buckets/{bucket_id}?force=1 if force
func (c *ActivityWatchClient) DeleteBucket(bucketID string, force bool) error {
	endpoint := fmt.Sprintf("buckets/%s", bucketID)
	if force {
		endpoint += "?force=1"
	}
	_, err := c.deleteReq(endpoint, nil)
	return err
}

// InsertEvent => POST /buckets/{bucket_id}/events
func (c *ActivityWatchClient) InsertEvent(bucketID string, evt *models.Event) error {
	endpoint := fmt.Sprintf("buckets/%s/events", bucketID)
	data := []interface{}{evt.ToJSONDict()}
	_, err := c.post(endpoint, data, nil)
	return err
}

// InsertEvents => multiple events
func (c *ActivityWatchClient) InsertEvents(bucketID string, evts []*models.Event) error {
	endpoint := fmt.Sprintf("buckets/%s/events", bucketID)
	var data []interface{}
	for _, e := range evts {
		data = append(data, e.ToJSONDict())
	}
	_, err := c.post(endpoint, data, nil)
	return err
}

// Heartbeat => POST /buckets/{bucket_id}/heartbeat?pulsetime=XYZ
func (c *ActivityWatchClient) Heartbeat(bucketID string, evt *models.Event, pulsetime float64, queued bool, commitInterval *float64) error {
	endpoint := fmt.Sprintf("buckets/%s/heartbeat?pulsetime=%f", bucketID, pulsetime)
	ci := c.CommitInterval
	if commitInterval != nil {
		ci = *commitInterval
	}

	if queued {
		last, ok := c.lastHeartbeat[bucketID]
		if !ok {
			c.lastHeartbeat[bucketID] = evt
			return nil
		}
		merged := utils.HeartbeatMerge(*last, *evt, pulsetime)
		if merged != nil {
			diff := merged.Duration.Seconds()
			if diff >= ci {
				data := merged.ToJSONDict()
				c.requestQueue.AddRequest(endpoint, data)
				c.lastHeartbeat[bucketID] = evt
			} else {
				c.lastHeartbeat[bucketID] = merged
			}
		} else {
			data := last.ToJSONDict()
			c.requestQueue.AddRequest(endpoint, data)
			c.lastHeartbeat[bucketID] = evt
		}
		return nil
	}
	// direct call
	_, err := c.post(endpoint, evt.ToJSONDict(), nil)
	return err
}

// Connect => start requestQueue
func (c *ActivityWatchClient) Connect() {
	c.requestQueue.Start()
}

// Disconnect => stop requestQueue
func (c *ActivityWatchClient) Disconnect() {
	c.requestQueue.Stop()
	c.requestQueue.Wait()

	// Recreate a fresh queue
	c.requestQueue = NewRequestQueue(c)
}

// WaitForStart => tries /info until success or timeout
func (c *ActivityWatchClient) WaitForStart(timeoutSeconds int) error {
	start := time.Now()
	sleepTime := 100 * time.Millisecond
	for {
		if time.Since(start).Seconds() > float64(timeoutSeconds) {
			return fmt.Errorf("server at %s did not start in time", c.ServerAddress)
		}
		_, err := c.GetInfo()
		if err == nil {
			break
		}
		time.Sleep(sleepTime)
		sleepTime *= 2
	}
	return nil
}

// RequestQueue => simulates a background thread for retrying requests
type RequestQueue struct {
	client           *ActivityWatchClient
	mu               sync.Mutex
	stopped          bool
	stopChan         chan struct{}
	wg               sync.WaitGroup
	connected        bool
	attemptReconnect time.Duration

	queueMu   sync.Mutex
	queue     []RequestQueueItem
	buckets   []BucketRegistration
	ticker    *time.Ticker
	tickerDur time.Duration
}

// NewRequestQueue => constructs the queue
func NewRequestQueue(client *ActivityWatchClient) *RequestQueue {
	return &RequestQueue{
		client:           client,
		stopChan:         make(chan struct{}),
		attemptReconnect: 10 * time.Second,
		queue:            []RequestQueueItem{},
		buckets:          []BucketRegistration{},
		tickerDur:        200 * time.Millisecond,
	}
}

// Start => spawns the goroutine
func (rq *RequestQueue) Start() {
	rq.ticker = time.NewTicker(rq.tickerDur)
	rq.wg.Add(1)
	go rq.run()
}

// Stop => signals goroutine
func (rq *RequestQueue) Stop() {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if rq.stopped {
		return
	}
	rq.stopped = true
	close(rq.stopChan)
}

// Wait => blocks until goroutine ends
func (rq *RequestQueue) Wait() {
	rq.wg.Wait()
}

func (rq *RequestQueue) run() {
	defer rq.wg.Done()
	for {
		select {
		case <-rq.stopChan:
			rq.client.logger.Println("RequestQueue stopping.")
			rq.ticker.Stop()
			return
		case <-rq.ticker.C:
			if !rq.connected {
				rq.tryConnect()
			} else {
				rq.dispatchRequest()
			}
		}
	}
}

func (rq *RequestQueue) tryConnect() {
	if err := rq.createBuckets(); err != nil {
		rq.client.logger.Printf("Not connected. Will retry in a bit: %v", err)
		rq.connected = false
	} else {
		rq.connected = true
		rq.client.logger.Printf("Connected. queue size=%d", len(rq.queue))
	}
}

func (rq *RequestQueue) createBuckets() error {
	for _, b := range rq.buckets {
		if err := rq.client.CreateBucket(b.ID, b.Type, false); err != nil {
			return err
		}
	}
	return nil
}

func (rq *RequestQueue) dispatchRequest() {
	rq.queueMu.Lock()
	defer rq.queueMu.Unlock()
	if len(rq.queue) == 0 {
		return
	}
	item := rq.queue[0]
	err := rq.dispatch(item)
	if err != nil {
		rq.connected = false
		rq.client.logger.Printf("Failed dispatch => %v", err)
		time.Sleep(500 * time.Millisecond)
		return
	}
	rq.queue = rq.queue[1:]
}

func (rq *RequestQueue) dispatch(item RequestQueueItem) error {
	_, err := rq.client.post(item.Endpoint, item.Data, nil)
	return err
}

func (rq *RequestQueue) AddRequest(endpoint string, data map[string]interface{}) {
	if endpoint == "" {
		return
	}
	rq.queueMu.Lock()
	defer rq.queueMu.Unlock()
	rq.queue = append(rq.queue, RequestQueueItem{Endpoint: endpoint, Data: data})
}

func (rq *RequestQueue) RegisterBucket(bucketID, eventType string) {
	rq.buckets = append(rq.buckets, BucketRegistration{ID: bucketID, Type: eventType})
}

func (rq *RequestQueue) IsAlive() bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	return !rq.stopped
}

func (c *ActivityWatchClient) IsQueueAlive() bool {
	return c.requestQueue.IsAlive()
}
