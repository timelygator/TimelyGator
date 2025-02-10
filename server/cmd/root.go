package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"timelygator/server/database"
	"timelygator/server/routes"
	"timelygator/server/utils/types"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tg-server",
	Short: "TimelyGator is a time tracking application. This is the server cli for the project",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := types.Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Error parsing config: %v", err)
		}
		if strings.Contains(cfg.DataSourceName, "@") {
			log.Fatalf("Only SQLite DSN is supported")
		}

		datastore, err := database.InitDB(cfg)
		if err != nil {
			log.Fatalf("Error initializing database: %v", err)
		}

		r := mux.NewRouter()
		routes.RegisterRoutes(cfg, datastore, r)

		slog.Info(fmt.Sprintf("Server running on :%s", cfg.Port))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
