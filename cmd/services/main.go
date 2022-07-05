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

	// Get config data
	cfgData, getCfgDataErr := config.Get().GetData(config.UPDATE_CONFIG_DATA)
	if getCfgDataErr != nil {
		log.Fatal("Error getting config data: ", getCfgDataErr)
	}

	// Initialize controller
	controllers.Initialize()

	// Create App
	rs := routes.CreateRoutingServer()

	// Start Server
	log.Print("Datastore service is ready...")

	addr := cfgData.HostName + cfgData.HostPort
	log.Print("The address used the service is: ", addr)
	log.Fatal(http.ListenAndServe(addr, rs.Router))
}
