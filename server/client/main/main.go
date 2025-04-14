package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/datatypes"
	"github.com/spf13/cobra"

	"timelygator/server/database/models"
	"timelygator/server/client"
)

// Global options that we parse from root flags
type rootOpts struct {
	host    string
	port    int
	verbose bool
	testing bool
}

// We'll store a single global TimelyGatorClient in our CLI. Another approach is to store this in the cobra command context.
var gClient *client.TimelyGatorClient
var gOpts = rootOpts{}

// rootCmd is the primary (root) CLI command. 
var rootCmd = &cobra.Command{
	Use:   "tg-cli",
	Short: "CLI utility for TimelyGatorClient to interact with server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// This hook is called before every subcommand. We can initialize the TimelyGatorClient here.

		// If user specified a port override
		finalPort := gOpts.port
		if gOpts.port == 8080 {
			if gOpts.testing {
				finalPort = 8080 // change
			} else {
				finalPort = 8080
			}
		}

		// Create the client
		gClient = client.NewTimelyGatorClient(
			"tg-cli",      // clientName
			gOpts.testing, // testing
			&gOpts.host,
			intToStringPtr(finalPort),
			"http", // or "http"
		)

		// If verbose => set logging to debug
		if gOpts.verbose {
			log.Printf("Running in verbose mode. Host=%s Port=%d\n", gOpts.host, finalPort)
		}
	},
}

// heartbeatCmd => `cli heartbeat <bucket_id> <data> [--pulsetime=60]`
var heartbeatCmd = &cobra.Command{
    Use:   "heartbeat <bucket_id> <data>",
    Short: "Send a heartbeat to bucket with JSON data",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        bucketID := args[0]
        dataStr := args[1]

        pulseTime, _ := strconv.Atoi(cmd.Flag("pulsetime").Value.String())

        // First, ensure dataStr is valid JSON
        if !json.Valid([]byte(dataStr)) {
            return fmt.Errorf("heartbeat error: provided data is not valid JSON:\n%v", dataStr)
        }

        // We store raw JSON as datatypes.JSON in the Event
        now := time.Now().UTC()
        evt := models.Event{
            Timestamp: now,
            Duration:  0,
            // Convert string to []byte for datatypes.JSON
            Data: datatypes.JSON([]byte(dataStr)),
        }

        log.Printf("Sending heartbeat: bucket=%s, data=%v, pulsetime=%d\n", bucketID, dataStr, pulseTime)
        if err := gClient.Heartbeat(bucketID, &evt, float64(pulseTime), false, nil); err != nil {
            return fmt.Errorf("heartbeat error: %v", err)
        }
        return nil
    },
}


// bucketsCmd => `tg-cli buckets`
var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "List all buckets",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := gClient.GetBucketsMap()
		if err != nil {
			return fmt.Errorf("failed to get buckets: %v", err)
		}
		log.Println("Buckets:")
		// The returned object is map[string]interface{}
		for key := range b {
			log.Printf(" - %s\n", key)
		}
		return nil
	},
}

// eventsCmd => `tg-cli events <bucket_id>`
var eventsCmd = &cobra.Command{
	Use:   "events <bucket_id>",
	Short: "Query events from the specified bucket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketID := args[0]
		// We'll just do "GetEvents" with no limit/time. Adjust later maybe
		evts, err := gClient.GetEvents(bucketID, -1, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to get events: %v", err)
		}
		log.Println("Events:")
		for _, e := range evts {
			// e is map[string]interface{} from the server
			tsRaw := e["timestamp"]
			durRaw := e["duration"]
			dataRaw := e["data"]
			log.Printf(" - TS=%v, Duration=%v, Data=%v\n", tsRaw, durRaw, dataRaw)
		}
		return nil
	},
}

// queryCmd => `tg-cli query <path> [--name ...] [--cache] [--json] [--start ...] [--stop ...]`
var queryCmd = &cobra.Command{
	Use:   "query <path>",
	Short: "Run a query in file at `path` on the server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		cacheFlag, _ := cmd.Flags().GetBool("cache")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		nameStr, _ := cmd.Flags().GetString("name")

		startStr, _ := cmd.Flags().GetString("start")
		stopStr, _ := cmd.Flags().GetString("stop")

		start, err := parseDateTime(startStr)
		if err != nil {
			return err
		}
		stop, err := parseDateTime(stopStr)
		if err != nil {
			return err
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}
		queryStr := string(bytes)

		timeperiods := [][2]time.Time{{start, stop}}
		res, err := gClient.Query(queryStr, timeperiods, &nameStr, cacheFlag)
		if err != nil {
			return fmt.Errorf("query error: %v", err)
		}

		if jsonFlag {
			out, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(out))
		} else {
			for i, period := range res {
				periodArr, ok := period.([]interface{})
				if !ok {
					continue
				}
				fmt.Printf("Showing 10 out of %d events in period #%d:\n", len(periodArr), i+1)
				for j, evtRaw := range periodArr {
					if j >= 10 {
						break
					}
					evtMap, _ := evtRaw.(map[string]interface{})
					// remove "id", "timestamp" for printing
					delete(evtMap, "id")
					delete(evtMap, "timestamp")
					dur := evtMap["duration"]
					dataVal := evtMap["data"]
					fmt.Printf(" - Duration: %v \tData: %v\n", dur, dataVal)
				}
			}
		}
		return nil
	},
}

// reportCmd => `tg-cli report <hostname> [--cache] [--start] [--stop] TODO [--limit=10]`
var reportCmd = &cobra.Command{
	Use:   "report <hostname>",
	Short: "Generate an activity report for a host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname := args[0]
		cacheFlag, _ := cmd.Flags().GetBool("cache")
		// limit, _ := cmd.Flags().GetInt("limit")

		startStr, _ := cmd.Flags().GetString("start")
		stopStr, _ := cmd.Flags().GetString("stop")

		start, err := parseDateTime(startStr)
		if err != nil {
			return err
		}
		stop, err := parseDateTime(stopStr)
		if err != nil {
			return err
		}

		// build query using "fullDesktopQuery" from queries or something
		bidWindow := fmt.Sprintf("tg-observer-window_%s", hostname)
		bidAfk := fmt.Sprintf("tg-observer-afk_%s", hostname)

		params := /* queries. */ client.DesktopQueryParams{
			QueryParams: client.QueryParams{
				BidBrowsers:    []string{},   // or fill with browser observer(s)
				Classes:        client.GetClasses(gClient),
				FilterClasses:  [][]string{}, // none
				FilterAfk:      true,
				IncludeAudible: true,
			},
			BidWindow: bidWindow,
			BidAfk:    bidAfk,
		}
		// build the query
		queryStr := client.FullDesktopQuery(gClient, &params)

		timeperiods := [][2]time.Time{{start, stop}}
		res, err := gClient.Query(queryStr, timeperiods, nil, cacheFlag)
		if err != nil {
			return fmt.Errorf("failed to query: %v", err)
		}
		if len(res) == 0 {
			fmt.Println("No data returned from server.")
			return nil
		}
		periodAny := res[0]
		periodMap, ok := periodAny.(map[string]interface{})
		if !ok {
			fmt.Println("No valid data in period.")
			return nil
		}

		// parse out categories, window, etc. 
		windowObj, _ := periodMap["window"].(map[string]interface{})
		// `ex`ample: cat_events => parse them
		catEventsRaw, _ := windowObj["cat_events"].([]interface{})
		// do a partial sum or display logic
		fmt.Printf("Category events. Found %d items.\n", len(catEventsRaw))

		return nil
	},
}

// exportCmd => `tg-cli export`
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all buckets and their events",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := gClient.ExportAll()
		if err != nil {
			return fmt.Errorf("failed to export all data: %v", err)
		}

		log.Println("Exported Data:")
		for bucketID, bucketData := range data {
			log.Printf(" - Bucket ID: %s\n", bucketID)
			if bucketInfo, ok := bucketData.(map[string]interface{}); ok {
				events, ok := bucketInfo["events"].([]interface{})
				if ok {
					log.Printf("   Events Count: %d\n", len(events))
				} else {
					log.Println("   No events found.")
				}
			}
		}

		return nil
	},
}


// canonicalCmd => `tg-cli canonical <hostname> [--cache] [--start] [--stop]`
var canonicalCmd = &cobra.Command{
	Use:   "canonical <hostname>",
	Short: "Query 'canonical events' for a single host (filtered, classified)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname := args[0]
		cacheFlag, _ := cmd.Flags().GetBool("cache")

		startStr, _ := cmd.Flags().GetString("start")
		stopStr, _ := cmd.Flags().GetString("stop")

		start, err := parseDateTime(startStr)
		if err != nil {
			return err
		}
		stop, err := parseDateTime(stopStr)
		if err != nil {
			return err
		}

		// we do something similar to CanonicalEvents( queries.DesktopQueryParams(...) )
		bidWindow := fmt.Sprintf("tg-observer-window_%s", hostname)
		bidAfk := fmt.Sprintf("tg-observer-afk_%s", hostname)

		params := &client.DesktopQueryParams{
			QueryParams: client.QueryParams{
				BidBrowsers:   []string{},
				Classes:       client.DefaultClasses, // or get from server
				FilterClasses: [][]string{},
				FilterAfk:     true,
			},
			BidWindow: bidWindow,
			BidAfk:    bidAfk,
		}

		qStr := client.CanonicalEvents(gClient, params)
		qStr += "\nRETURN = events;"
		timeperiods := [][2]time.Time{{start, stop}}
		res, err := gClient.Query(qStr, timeperiods, nil, cacheFlag)
		if err != nil {
			return fmt.Errorf("failed to run canonical query: %v", err)
		}

		for _, periodAny := range res {
			periodSlice, ok := periodAny.([]interface{})
			if !ok {
				continue
			}
			fmt.Printf("Showing last 10 out of %d events:\n", len(periodSlice))
			for i := len(periodSlice) - 10; i < len(periodSlice); i++ {
				if i < 0 {
					continue
				}
				evtRaw := periodSlice[i]
				evtMap, _ := evtRaw.(map[string]interface{})
				tsVal := evtMap["timestamp"]
				durVal := evtMap["duration"]
				dataVal, _ := evtMap["data"].(map[string]interface{})
				title := ""
				if t, ok := dataVal["title"].(string); ok {
					title = t
				}
				// shorten the title to 60?
				if len(title) > 60 {
					title = title[:60] + "..."
				}
				appStr := ""
				if a, ok := dataVal["app"].(string); ok {
					appStr = "[" + a + "] "
				}
				fmt.Printf(" %v \t%v \t%v%v\n", tsVal, durVal, appStr, title)
			}
			// sum durations
			totalDur := 0.0
			for _, e := range periodSlice {
				evtMap, _ := e.(map[string]interface{})
				dVal, _ := evtMap["duration"].(float64)
				totalDur += dVal
			}
			fmt.Printf("\nTotal duration: %f seconds\n", totalDur)
		}
		return nil
	},
}

// parseDateTime => a minimal date/time parser that tries "RFC3339" or fallback
func parseDateTime(s string) (time.Time, error) {
	if s == "" {
		// default to time.Now()
		return time.Now().UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}
	// fallback parse? or return error
	return time.Time{}, fmt.Errorf("invalid time format (must be RFC3339), got: %q", s)
}

// intToStringPtr => quick helper to turn an int into a *string
func intToStringPtr(i int) *string {
	s := strconv.Itoa(i)
	return &s
}

// init initializes root flags and subcommands
func init() {
	rootCmd.PersistentFlags().StringVar(&gOpts.host, "host", "localhost", "Address of host")
	rootCmd.PersistentFlags().IntVar(&gOpts.port, "port", 8080, "Port to use")
	rootCmd.PersistentFlags().BoolVar(&gOpts.testing, "testing", false, "Use testing mode (port=8080 if not specified)") // change
	rootCmd.PersistentFlags().BoolVar(&gOpts.verbose, "verbose", false, "Enable verbose logging")

	// Subcommand: heartbeat
	heartbeatCmd.Flags().Int("pulsetime", 60, "Pulsetime for merging heartbeats")

	// Subcommand: query
	queryCmd.Flags().String("name", "", "Optional query name (for caching)")
	queryCmd.Flags().Bool("cache", false, "Use caching")
	queryCmd.Flags().Bool("json", false, "Output JSON")
	queryCmd.Flags().String("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339), "Start time (RFC3339)")
	queryCmd.Flags().String("stop", time.Now().Add(365*24*time.Hour).Format(time.RFC3339), "Stop time (RFC3339)")

	// Subcommand: report
	reportCmd.Flags().Bool("cache", false, "Use caching")
	reportCmd.Flags().String("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339), "Start time (RFC3339)")
	reportCmd.Flags().String("stop", time.Now().Add(365*24*time.Hour).Format(time.RFC3339), "Stop time (RFC3339)")
	// reportCmd.Flags().Int("limit", 10, "Limit for top items")

	// Subcommand: canonical
	canonicalCmd.Flags().Bool("cache", false, "Use caching")
	canonicalCmd.Flags().String("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339), "Start time (RFC3339)")
	canonicalCmd.Flags().String("stop", time.Now().Add(365*24*time.Hour).Format(time.RFC3339), "Stop time (RFC3339)")

	// Register subcommands
	rootCmd.AddCommand(heartbeatCmd)
	rootCmd.AddCommand(bucketsCmd)
	rootCmd.AddCommand(eventsCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(canonicalCmd)
}

// main calls Execute on rootCmd
func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
