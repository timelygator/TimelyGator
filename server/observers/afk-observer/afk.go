package afkobserver

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"

    "gorm.io/datatypes"

    "timelygator/server/database/models"
    "timelygator/server/client"
)

type AFKWatcher struct {
    client     *client.TimelyGatorClient
    bucketName string

    timeout   time.Duration
    pollTime  time.Duration
    isTesting bool
    verbose   bool
}

func NewAFKWatcher(timeout, pollTime time.Duration, host, port string, testing, verbose bool) *AFKWatcher {
    tgClient := client.NewTimelyGatorClient(
        "tg-observer-afk",
        testing,
        &host,
        &port,
        "http",
    )

    bucketName := fmt.Sprintf("%s_%s", tgClient.ClientName, tgClient.ClientHostname)

    return &AFKWatcher{
        client:     tgClient,
        bucketName: bucketName,
        timeout:    timeout,
        pollTime:   pollTime,
        isTesting:  testing,
        verbose:    verbose,
    }
}

// Run starts the AFKWatcher, including signal handling and the heartbeat loop
func (w *AFKWatcher) Run() {
    if runtime.GOOS != "linux" {
        log.Fatalf("Unsupported platform: %s (only Linux implemented)", runtime.GOOS)
    }

    log.Println("[AFKWatcher] started")

    if err := w.client.WaitForStart(10); err != nil {
        log.Fatalf("Failed to start client: %v", err)
        return
    }

    log.Println("[AFKWatcher] Creating bucket:", w.bucketName)
    if err := w.client.CreateBucket(w.bucketName, "afkstatus", false); err != nil {
        log.Fatalf("Failed to create bucket: %v", err)
        return
    }

    go w.heartbeatLoop()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    <-sigChan
    log.Println("[AFKWatcher] stopped by signal")
}

// heartbeatLoop regularly checks AFK status and sends heartbeats to the client
func (w *AFKWatcher) heartbeatLoop() {
    var afk bool

    ticker := time.NewTicker(w.pollTime)
    defer ticker.Stop()
    for {
        // Stop if the parent process has died
        if os.Getppid() == 1 {
            log.Println("[AFKWatcher] stopped because parent process died")
            return
        }

        now := time.Now().UTC()
        secondsSinceInput := SecondsSinceLastInput()
        // compute the last‐input time by subtracting secondsSinceInput
        lastInput := now.Add(-time.Duration(secondsSinceInput * float64(time.Second)))
        if w.verbose {
            log.Printf("Seconds since last input: %.2f\n", secondsSinceInput)
        }

        // Handle AFK state transitions
        if afk && secondsSinceInput < w.timeout.Seconds() {
            log.Println("[AFKWatcher] No longer AFK")
            w.ping(afk, lastInput, 0)
            afk = false

            // Ensure the latest event is not missed
            w.ping(afk, lastInput.Add(1*time.Millisecond), 0)

        } else if !afk && secondsSinceInput >= w.timeout.Seconds() {
            log.Println("[AFKWatcher] Became AFK")
            w.ping(afk, lastInput, 0)
            afk = true

            w.ping(afk, lastInput.Add(1*time.Millisecond), secondsSinceInput)

        } else {
            // Send a regular heartbeat if state has not changed
            if afk {
                w.ping(afk, lastInput.Add(1*time.Millisecond), secondsSinceInput)
            } else {
                w.ping(afk, lastInput, 0)
            }
        }

        <-ticker.C
    }
}

// ping sends a heartbeat to the TimelyGatorClient
func (w *AFKWatcher) ping(afk bool, timestamp time.Time, durationSeconds float64) {
    status := "not-afk"
    if afk {
        status = "afk"
    }
    // Prepare event data as JSON
    dataMap := map[string]interface{}{
        "status": status,
    }
    rawJSON, _ := json.Marshal(dataMap) // Convert map to JSON (ignoring errors for brevity)

    ev := &models.Event{
        Timestamp: timestamp,
        Duration:  durationSeconds,
        Data:      datatypes.JSON(rawJSON),
    }

    pulsetime := w.timeout.Seconds() + w.pollTime.Seconds()
    interval := float64(60.0)

    // Send the heartbeat to the TimelyGatorClient
    w.client.Heartbeat(w.bucketName, ev, pulsetime, false, &interval)
}
