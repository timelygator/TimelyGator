package main

import (
	"log/slog"
	"os"
	"timelygator/server/cmd"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Try to load .env file
	godotenv.Load()

	cmd.Execute()
}
