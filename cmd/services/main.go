package main

import (
	"log"
	"net/http"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/controllers"
	"github.com/sflewis2970/datastore-service/routes"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Create config object and load config data into memory
	log.Print("Loading config data...")
	cfg := config.GetConfig()

	// Load config
	controllers.LoadConfig()

	// Create App
	rs := routes.CreateRoutingServer()

	// Start Server
	log.Print("Datastore service is ready...")

	addr := cfg.GetConfigData().HostName + cfg.GetConfigData().HostPort
	log.Print("The address used the service is: ", addr)
	log.Fatal(http.ListenAndServe(addr, rs.Router))
}
