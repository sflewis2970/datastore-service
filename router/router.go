package router

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/sflewis2970/datastore-service/controllers"
)

type MessageRouter struct {
	MuxRouter *mux.Router
}

var msgRouter *MessageRouter

func (rs *MessageRouter) setupRoutes() {
	// Display log message
	log.Print("Setting up Datastore service routes")

	// Setup routes
	rs.MuxRouter.HandleFunc("/api/v1/ds/status", controllers.Status).Methods("GET")
	rs.MuxRouter.HandleFunc("/api/v1/ds/insert", controllers.Insert).Methods("POST")
	rs.MuxRouter.HandleFunc("/api/v1/ds/get", controllers.Get).Methods("POST")
}

func New() *MessageRouter {
	msgRouter = new(MessageRouter)

	// Create router
	msgRouter.MuxRouter = mux.NewRouter()

	// Setting up routes
	msgRouter.setupRoutes()

	return msgRouter
}
