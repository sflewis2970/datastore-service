package routes

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/sflewis2970/datastore-service/controllers"
)

type RoutingServer struct {
	Router *mux.Router
}

func (rs *RoutingServer) setupRoutes() {
	// Display log message
	log.Print("Setting up Datastore service routes")

	// Setup routes
	rs.Router.HandleFunc("/", controllers.Home).Methods("GET")
	rs.Router.HandleFunc("/api/v1/ds/status", controllers.Status).Methods("GET")
	rs.Router.HandleFunc("/api/v1/ds/insert", controllers.Insert).Methods("POST")
	rs.Router.HandleFunc("/api/v1/ds/get", controllers.Get).Methods("POST")
}

func CreateRoutingServer() *RoutingServer {
	rs := new(RoutingServer)

	// Create router
	rs.Router = mux.NewRouter()

	// Setting up routes
	rs.setupRoutes()

	return rs
}
