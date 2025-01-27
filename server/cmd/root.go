package cmd

import (
	"fmt"
	"log"
	"net/http"
	"timelygator/server/database"
	"timelygator/server/routes"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "Time Tracker Application",
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.InitDB(); err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}

		r := mux.NewRouter()
		routes.RegisterRoutes(r)

		fmt.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
