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

	// Get config object
	cfg, getErr := config.GetConfig()
	if getErr != nil {
		log.Fatal("Error getting config: ", getErr)
	}

	// Load config data
	cfgData := cfg.GetConfigData()

	// Initialize controller
	controllers.InitializeController()

	// Create App
	rs := routes.CreateRoutingServer()

	// Start Server
	log.Print("Datastore service is ready...")

	addr := cfgData.HostName + cfgData.HostPort
	log.Print("The address used the service is: ", addr)
	log.Fatal(http.ListenAndServe(addr, rs.Router))
}
