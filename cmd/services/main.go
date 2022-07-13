package main

import (
	"log"
	"net/http"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/controllers"
	"github.com/sflewis2970/datastore-service/router"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get config data
	cfgData, cfgDataErr := config.Get().GetData(config.REFRESH_CONFIG_DATA)
	if cfgDataErr != nil {
		log.Fatal("Error getting config data: ", cfgDataErr)
	}

	// Initialize controller
	controllers.New()

	// Create App
	msgRouter := router.New()

	// Start Server
	log.Print("Datastore service is ready...")

	log.Print("Host: ", cfgData.Host)
	log.Print("Port: ", cfgData.Port)
	addr := cfgData.Host + ":" + cfgData.Port
	log.Print("The address used the service is: ", addr)
	log.Fatal(http.ListenAndServe(addr, msgRouter.MuxRouter))
}
