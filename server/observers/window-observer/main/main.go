//go:build linux || darwin || windows

// Command window-observer is the TimelyGator port of aw‑watcher‑window
// written in Go.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"timelygator/server/client"
	"timelygator/server/database/models"
	"timelygator/server/observers/window-observer"
	"timelygator/server/observers/window-observer/lib"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// ----- defaults from environment -------------------------------
	cfg, err := windowobserver.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// ----- CLI flags ------------------------------------------------
	var (
		host          = flag.String("host", cfg.Host, "Server host")
		port          = flag.String("port", cfg.Port, "Server port")
		testing       = flag.Bool("testing", cfg.Testing, "Testing mode")
		verbose       = flag.Bool("verbose", cfg.Verbose, "Verbose logging")
		excludeTitle  = flag.Bool("exclude-title", cfg.ExcludeTitle, "Replace every title with 'excluded'")
		excludeTitles = flag.String("exclude-titles", strings.Join(cfg.ExcludeTitles, ","), "Comma‑separated regex list to anonymize titles")
		pollTime      = flag.Float64("poll-time", cfg.PollTime, "Polling interval in seconds")
		strategy      = flag.String("strategy", cfg.Strategy, "macOS only: jxa | applescript | swift")
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	// ----- macOS permissions ---------------------------------------
	if runtime.GOOS == "darwin" {
		windowobserver.BackgroundEnsurePermissions()
	}

	// ----- TimelyGator client --------------------------------------
	tg := client.NewTimelyGatorClient("tg-observer-window", *testing, host, port, "http")
	if err := tg.WaitForStart(10); err != nil {
		log.Fatalf("server not ready: %v", err)
	}

	bucketID := fmt.Sprintf("%s_%s", tg.ClientName, tg.ClientHostname)
	if err := tg.CreateBucket(bucketID, "currentwindow", false); err != nil {
		log.Fatalf("create bucket: %v", err)
	}

	// ----- swift shortcut (macOS) ----------------------------------
	if runtime.GOOS == "darwin" && *strategy == "swift" {
		bin := filepath.Join(filepath.Dir(os.Args[0]), "aw-watcher-window-macos")
		cmd := exec.Command(bin, tg.ServerAddress, bucketID, tg.ClientHostname, tg.ClientName)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("failed to start swift helper: %v", err)
		}
		// clean up on exit
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
			<-sig
			_ = cmd.Process.Kill()
		}()
		_ = cmd.Wait()
		return
	}

	// compile exclude regexes
	var patterns []*regexp.Regexp
	if *excludeTitles != "" {
		for _, raw := range strings.Split(*excludeTitles, ",") {
			if r, err := regexp.Compile("(?i)" + strings.TrimSpace(raw)); err == nil {
				patterns = append(patterns, r)
			} else {
				log.Printf("bad regex %q: %v", raw, err)
			}
		}
	}

	pollDur := time.Duration(*pollTime * float64(time.Second))
	ticker := time.NewTicker(pollDur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			handleHeartBeat(tg, bucketID, *strategy, *excludeTitle, patterns, pollDur)
		case <-interrupt():
			log.Println("window-observer shutting down")
			return
		}
	}
}

func handleHeartBeat(tg *client.TimelyGatorClient, bucket string, strategy string, exclTitle bool, patterns []*regexp.Regexp, poll time.Duration) {
	win, err := lib.GetCurrentWindow(strategy)
	if err != nil {
		log.Printf("getCurrentWindow error: %v", err)
		return
	}
	if win == nil {
		return
	}

	// anonymize titles if requested
	if exclTitle {
		win.Title = "excluded"
	} else {
		for _, r := range patterns {
			if r.MatchString(win.Title) {
				win.Title = "excluded"
				break
			}
		}
	}

	data := map[string]interface{}{"app": win.App, "title": win.Title}
	if win.URL != nil {
		data["url"] = *win.URL
	}
	if win.Incognito != nil {
		data["incognito"] = *win.Incognito
	}

	raw, _ := json.Marshal(data)
	ev := &models.Event{
		Timestamp: time.Now().UTC(),
		Duration:  0,
		Data:      raw,
	}

	pulse := poll.Seconds() + 1.0
	_ = tg.Heartbeat(bucket, ev, pulse, false, nil)
}

// interrupt returns a channel closed on SIGINT/SIGTERM.
func interrupt() <-chan struct{} {
	ch := make(chan struct{})
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc
		close(ch)
	}()
	return ch
}

func hostPtr(s *string) string { return *s }
func portPtr(i *int) string    { return fmt.Sprintf("%d", *i) }
