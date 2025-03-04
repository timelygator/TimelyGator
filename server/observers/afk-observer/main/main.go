package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"timelygator/server/observers/afk-observer"

	"github.com/spf13/cobra"
)

type Args struct {
	Host     string
	Port     string
	Testing  bool
	Verbose  bool
	Timeout  time.Duration
	PollTime time.Duration
}

func main() {
	cfg, err := afkobserver.LoadAFKConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var args Args

	var rootCmd = &cobra.Command{
		Use:   "afk-observer",
		Short: "A observer for keyboard/mouse input to detect AFK state",
		Run: func(cmd *cobra.Command, cmdArgs []string) {
			fmt.Println("=== Final Arguments ===")
			fmt.Printf("Host:     %s\n", args.Host)
			fmt.Printf("Port:     %s\n", args.Port)
			fmt.Printf("Testing:  %v\n", args.Testing)
			fmt.Printf("Verbose:  %v\n", args.Verbose)
			fmt.Printf("Timeout:  %v\n", args.Timeout)
			fmt.Printf("PollTime: %v\n", args.PollTime)

			watcher := afkobserver.NewAFKWatcher(args.Timeout, args.PollTime, args.Host, args.Port, args.Testing, args.Verbose)
			watcher.Run()
		},
	}

	rootCmd.Flags().StringVar(
		&args.Host, "host", "localhost", "Host (interface) for server",
	)
	rootCmd.Flags().StringVar(
		&args.Port, "port", "8080", "Port for server",
	)
	rootCmd.Flags().BoolVar(
		&args.Testing, "testing", false, "Run in testing mode",
	)
	rootCmd.Flags().BoolVar(
		&args.Verbose, "verbose", false, "Enable verbose logging",
	)
	rootCmd.Flags().DurationVar(
		&args.Timeout, "timeout", time.Duration(cfg.Timeout)*time.Second, "AFK timeout (e.g. 180s, 5m, 1h)",
	)
	rootCmd.Flags().DurationVar(
		&args.PollTime, "poll-time", time.Duration(cfg.PollTime)*time.Second, "Poll time (e.g. 5s, 1m)",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
