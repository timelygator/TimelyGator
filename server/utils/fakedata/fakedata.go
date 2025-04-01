package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "strings"
    "time"

    "github.com/spf13/cobra"
    "gorm.io/datatypes"

    "timelygator/server/client"
    "timelygator/server/database/models"
)

// sampleData holds the “template” for an event (e.g., window, afk, browser).
// Weight is used for probability. Minutes is a “typical” event duration.
type sampleData struct {
    App     string
    Title   string
    Status  string
    URL     string
    Weight  int     // Probability weight
    Minutes float64 // Base duration for the event in minutes
}

var (
    sampleDataAFK = []sampleData{
        {Status: "not-afk", Weight: 1, Minutes: 120},
        {Status: "afk", Weight: 1, Minutes: 10},
    }

    sampleDataWindow = []sampleData{
        // Meetings
        {App: "zoom", Title: "Zoom Meeting", Weight: 3, Minutes: 20},
        // Games
        {App: "Minecraft", Title: "Minecraft", Weight: 2, Minutes: 200},
        // timelygator-related
        {App: "Firefox", Title: "TimelyGator/timelygator: Track how you spend your time", Weight: 20, Minutes: 5},
        {App: "Terminal", Title: "vim ~/code/timelygator/other/tg-fakedata", Weight: 10},
        {App: "Terminal", Title: "vim ~/code/timelygator/README.md", Weight: 3, Minutes: 5},
        {App: "Terminal", Title: "vim ~/code/timelygator/tg-server", Weight: 5},
        {App: "Terminal", Title: "bash ~/code/timelygator", Weight: 5},
        // Misc work
        {App: "Firefox", Title: "Gmail - mail.google.com/", Weight: 5, Minutes: 10},
        {App: "Firefox", Title: "Stack Overflow - stackoverflow.com/", Weight: 10, Minutes: 5},
        {App: "Firefox", Title: "Google Calendar - calendar.google.com/", Weight: 5, Minutes: 2},
        // Social media
        {App: "Firefox", Title: "reddit: the front page of the internet - reddit.com/", Weight: 10, Minutes: 10},
        {App: "Firefox", Title: "Home / Twitter - twitter.com/", Weight: 10, Minutes: 8},
        {App: "Firefox", Title: "Facebook - facebook.com/", Weight: 10, Minutes: 3},
        {App: "Chrome", Title: "Unknown site", Weight: 2},
        // Media
        {App: "Spotify", Title: "Spotify", Weight: 8, Minutes: 3},
        {App: "Chrome", Title: "YouTube - youtube.com/", Weight: 4, Minutes: 25},
    }

    sampleDataBrowser = []sampleData{
        {Title: "GitHub", URL: "https://github.com", Weight: 10, Minutes: 10},
        {Title: "Twitter", URL: "https://twitter.com", Weight: 3, Minutes: 5},
        {Title: "YouTube", URL: "https://youtube.com", Weight: 5, Minutes: 20},
    }
)

const (
    hostname            = "fakedata"
    clientName          = "tg-fakedata"
    bucketWindow        = "tg-observer-window_" + hostname
    bucketAFK           = "tg-observer-afk_" + hostname
    bucketBrowserChrome = "tg-observer-web-chrome_" + hostname
    bucketBrowserFF     = "tg-observer-web-firefox_" + hostname
)

// Flags
var (
    sinceFlag string
    untilFlag string
)

// RootCmd for Cobra
var rootCmd = &cobra.Command{
    Use:   "tg-fakedata",
    Short: "Generate fake data for TimelyGator",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runFakeData()
    },
}

func main() {
    rootCmd.Flags().StringVar(&sinceFlag, "since", "",
        "Start date (YYYY-MM-DD). Defaults to 14 days before --until if omitted.")
    rootCmd.Flags().StringVar(&untilFlag, "until", "",
        "End date (YYYY-MM-DD). Defaults to today if omitted.")

    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func runFakeData() error {
    now := time.Now().UTC()

    // Parse or default the until date
    var until time.Time
    if untilFlag == "" {
        until = now
    } else {
        t, err := parseDateFlag(untilFlag)
        if err != nil {
            return fmt.Errorf("failed to parse --until date: %w", err)
        }
        until = t.UTC()
    }

    // Parse or default the since date
    var since time.Time
    if sinceFlag == "" {
        since = until.AddDate(0, 0, -14) // default to 14 days prior
    } else {
        t, err := parseDateFlag(sinceFlag)
        if err != nil {
            return fmt.Errorf("failed to parse --since date: %w", err)
        }
        since = t.UTC()
    }

    fmt.Printf("Range: %s to %s\n", since, until)

    // Create the client
    emptyString := ""
    c := client.NewTimelyGatorClient(clientName, false, &emptyString, &emptyString, emptyString)

    if err := c.CreateBucket(bucketWindow, "currentwindow", false); err != nil {
        return fmt.Errorf("failed to create window bucket: %w", err)
    }
    if err := c.CreateBucket(bucketAFK, "afkstatus", false); err != nil {
        return fmt.Errorf("failed to create AFK bucket: %w", err)
    }
    if err := c.CreateBucket(bucketBrowserChrome, "web.tab.current", false); err != nil {
        return fmt.Errorf("failed to create Chrome bucket: %w", err)
    }
    if err := c.CreateBucket(bucketBrowserFF, "web.tab.current", false); err != nil {
        return fmt.Errorf("failed to create Firefox bucket: %w", err)
    }

    // 2) Generate fake data
    buckets := generateAllDays(since, until)

    // 3) Insert events into DB or server
    for bucketID, evts := range buckets {
        for i := 0; i < len(evts); i++ {
            evts[i].BucketID = bucketID
        }

        // Log events before inserting
        for _, evt := range evts {
            log.Printf("%+v\n", evt)
        }

        // Convert evts to a slice of interface{} for your InsertEvents signature
        eventsInterface := make([]interface{}, len(evts))
        for i, evt := range evts {
            eventsInterface[i] = evt
        }

        if err := c.InsertEvents(bucketID, eventsInterface); err != nil {
            return fmt.Errorf("failed to insert events to bucket %q: %w", bucketID, err)
        }
        fmt.Printf("Inserted %d events into bucket %s\n", len(evts), bucketID)
    }

    return nil
}

func parseDateFlag(val string) (time.Time, error) {
    layout := "2006-01-02"
    return time.Parse(layout, val)
}

// generateAllDays iterates from start to end, day by day.
func generateAllDays(start, end time.Time) map[string][]models.Event {
    rand.Seed(int64(start.Unix() + end.Unix())) // So consistent for same range

    results := make(map[string][]models.Event)
    for d := start; d.Before(end) || sameDay(d, end); d = d.AddDate(0, 0, 1) {
        dayEvents := generateDay(d, end)
        for bucketID, evts := range dayEvents {
            results[bucketID] = append(results[bucketID], evts...)
        }
    }
    return results
}

// sameDay checks if two times share the same calendar date (UTC).
func sameDay(t1, t2 time.Time) bool {
    y1, m1, d1 := t1.UTC().Date()
    y2, m2, d2 := t2.UTC().Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}

// generateDay picks a random start time (08:00), random day length, and then splits in half for a “break”.
func generateDay(day, globalEnd time.Time) map[string][]models.Event {
    res := make(map[string][]models.Event)

    start := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, time.UTC)
    if start.After(globalEnd) {
        return res
    }

    isWeekday := start.Weekday() >= time.Monday && start.Weekday() <= time.Friday
    var dayDuration time.Duration
    if isWeekday {
        // 5 to 10 hours
        dayDuration = time.Duration(float64(time.Hour)*5 + rand.Float64()*float64(time.Hour)*5)
    } else {
        // 1 to 5 hours
        dayDuration = time.Duration(float64(time.Hour)*1 + rand.Float64()*float64(time.Hour)*4)
    }

    stop := start.Add(dayDuration)
    if stop.After(globalEnd) {
        stop = globalEnd
    }

    // Break in the middle
    breakStart := start.Add((stop.Sub(start)) / 2)
    breakDuration := time.Duration(60+rand.Intn(60)) * time.Minute // 60-120
    breakStop := breakStart.Add(breakDuration)
    if breakStop.After(stop) {
        breakStop = stop
    }

    // Activity from [start, breakStart] and [breakStop, stop + breakDuration]
    act1 := generateActivity(start, breakStart)
    act2 := generateActivity(breakStop, stop.Add(breakDuration))

    for k, v := range act1 {
        res[k] = append(res[k], v...)
    }
    for k, v := range act2 {
        res[k] = append(res[k], v...)
    }

    return res
}

// generateActivity builds AFK events, window events, browser events.
func generateActivity(start, end time.Time) map[string][]models.Event {
    res := make(map[string][]models.Event)
    if end.Before(start) {
        return res
    }

    // Generate AFK
    afkEvents := randomEvents(start, end, sampleDataAFK, 120*60)
    var windowEvents []models.Event
    var chromeEvents []models.Event
    var firefoxEvents []models.Event

    // For each non-afk chunk, generate window events
    for _, e := range afkEvents {
        status := getString(e.Data, "status")
        if status == "not-afk" {
            ts := e.Timestamp
            evEnd := ts.Add(e.Duration)
            evs := randomEvents(ts, evEnd, sampleDataWindow, 120*60)
            windowEvents = append(windowEvents, evs...)
        }
    }

    // For each window event that is Chrome/Firefox, generate tab events
    for _, w := range windowEvents {
        app := strings.ToLower(getString(w.Data, "app"))
        evEnd := w.Timestamp.Add(w.Duration)
        switch app {
        case "chrome":
            chromeEvents = append(chromeEvents, randomEvents(w.Timestamp, evEnd, sampleDataBrowser, 120*60)...)
        case "firefox":
            firefoxEvents = append(firefoxEvents, randomEvents(w.Timestamp, evEnd, sampleDataBrowser, 120*60)...)
        }
    }

    res[bucketAFK] = afkEvents
    res[bucketWindow] = windowEvents
    res[bucketBrowserChrome] = chromeEvents
    res[bucketBrowserFF] = firefoxEvents

    return res
}

// randomEvents uses weightedChoice on the given sample data to build consecutive events up to [stop].
func randomEvents(start, stop time.Time, samples []sampleData, maxSecs float64) []models.Event {
    var results []models.Event
    ts := start

    for ts.Before(stop) {
        s := weightedChoice(samples)
        d := pickDuration(s.Minutes, maxSecs)
        evEnd := ts.Add(d)
        if evEnd.After(stop) {
            evEnd = stop
        }
        dur := evEnd.Sub(ts)
        if dur <= 0 {
            break
        }

        // Build the data as a map, then encode to JSON
        dataMap := map[string]interface{}{}
        if s.App != "" {
            dataMap["app"] = s.App
        }
        if s.Title != "" {
            dataMap["title"] = s.Title
        }
        if s.Status != "" {
            dataMap["status"] = s.Status
        }
        if s.URL != "" {
            dataMap["url"] = s.URL
        }
        jsonData, _ := json.Marshal(dataMap)

        e := models.Event{
            Timestamp: ts,
            Duration:  dur,
            Data:      datatypes.JSON(jsonData),
        }
        results = append(results, e)

        ts = evEnd
    }
    return results
}

func pickDuration(minutes float64, maxSecs float64) time.Duration {
    if minutes > 0 {
        // random factor from 0.5..2.0
        f := 0.5 + rand.Float64()*1.5
        return time.Duration(f * minutes * float64(time.Minute))
    }
    // no base duration means pick random from 5 seconds..maxSecs
    sec := 5.0 + rand.Float64()*(maxSecs-5.0)
    return time.Duration(sec * float64(time.Second))
}

func weightedChoice(items []sampleData) sampleData {
    var total int
    for _, it := range items {
        total += it.Weight
    }
    r := rand.Intn(total)
    for _, it := range items {
        if r < it.Weight {
            return it
        }
        r -= it.Weight
    }
    return items[len(items)-1]
}

func getString(js datatypes.JSON, key string) string {
    // Quick helper: decode the JSON into a map to retrieve a field.
    var m map[string]interface{}
    _ = json.Unmarshal(js, &m)
    val, _ := m[key].(string)
    return val
}