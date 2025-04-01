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

