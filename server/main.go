package main

import (
	"log"
	"timelygator/server/cmd"

	"github.com/joho/godotenv"
)

//	@title			TimelyGator Server API
//	@version		0.1
//	@description	This is the API documentation for the TimelyGator Server API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://github.com/timelygator/timelygator

//	@license.name	GPLv3
//	@license.url	https://github.com/timelygator/TimelyGator/blob/main/LICENSE

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
