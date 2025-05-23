package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"timelygator/server/api"
	"timelygator/server/database"
	"timelygator/server/utils/types"

	"github.com/rs/cors"

	_ "timelygator/server/docs" // docs is generated by Swaggo

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
)

var rootCmd = &cobra.Command{
	Use:     "tg-server",
	Short:   "TimelyGator is a time tracking application. This is the server cli for the project",
	Version: types.ModuleVersion,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := types.Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Error parsing config: %v", err)
		}
		if strings.Contains(cfg.DataSourceName, "@") {
			log.Fatalf("Only SQLite DSN is supported")
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		})

		datastore, err := database.InitDB(cfg)
		if err != nil {
			log.Fatalf("Error initializing database: %v", err)
		}
		routes := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		api.RegisterRoutes(cfg, datastore, routes)

		router := mux.NewRouter()
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		router.PathPrefix("/api/v1/").Handler(routes)

		handler := c.Handler(router)

		slog.Info(fmt.Sprintf("Server running on %s:%s", cfg.Interface, cfg.Port))
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Interface, cfg.Port), handler))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
