package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/controllers"
	"github.com/sflewis2970/datastore-service/routes"
)

func setCfgEnv() {
	// Set hostname environment variable
	os.Setenv(config.HOSTNAME, "")

	// Set hostport environment variable
	os.Setenv(config.HOSTPORT, ":9090")

	// Set activedriver environment variable
	os.Setenv(config.ACTIVEDRIVER, "go-cache")

	// Set Go-cache environment variable
	switch os.Getenv(config.ACTIVEDRIVER) {
	case "go-cache":
		os.Setenv(config.DEFAULT_EXPIRATION, "1")
		os.Setenv(config.CLEANUP_INTERVAL, "30")
	case "mysql":
		os.Setenv(config.MYSQL_CONNECTION, "root:devStation@tcp(127.0.0.1:3306)/")
	case "postgres":
		os.Setenv(config.POSTGRES_HOST, "127.0.0.1")
		os.Setenv(config.POSTGRES_PORT, "5432")
		os.Setenv(config.POSTGRES_USER, "postgres")
	}

	// Response messages
	os.Setenv(config.CONGRATS, "Congrats! That is correct")
	os.Setenv(config.TRYAGAIN, "Nice try! Better luck on the next question")
}

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	useCfgFile := os.Getenv("USECONFIGFILE")
	if len(useCfgFile) == 0 {
		setCfgEnv()
	}

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
