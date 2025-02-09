package cmd

import (
	"fmt"
	"log"
	"net/http"
	"timelygator/server/database"
	"timelygator/server/routes"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type Config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"development"`
	Domain             string `env:"DOMAIN" envDefault:"localhost"`
	Port               string `env:"PORT" envDefault:"8080"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

var rootCmd = &cobra.Command{
	Use:   "tg-server",
	Short: "TimelyGator is a time tracking application. This is the server cli for the project",
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.InitDB(); err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}

		cfg := Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Could not parse environment variables: %v", err)
		}

		r := mux.NewRouter()
		routes.RegisterRoutes(r)

		fmt.Printf("Server running on port :%s\n", cfg.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
